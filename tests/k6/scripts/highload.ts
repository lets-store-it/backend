// @ts-ignore
import { check, sleep } from 'k6';
// @ts-ignore
import http from 'k6/http';
// @ts-ignore
import { SharedArray } from 'k6/data';
// @ts-ignore
import { randomString, randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

declare global {
  const __ENV: {
    BASE_URL?: string;
    TEST_USERS_COUNT?: string;
    VUS?: string;
    DURATION?: string;
  };
}

const BASE_URL = __ENV.BASE_URL || 'http://host.docker.internal:8080';
const TEST_USERS_COUNT = parseInt(__ENV.TEST_USERS_COUNT || '10');
const VUS = parseInt(__ENV.VUS || '5');
const DURATION = __ENV.DURATION || '1m';
const MAX_RETRIES = 3;
const RETRY_INTERVAL = 5;

interface TestUser {
  email: string;
  firstName: string;
  lastName: string;
}

const users = new SharedArray<TestUser>('users', function () {
  return Array.from({ length: TEST_USERS_COUNT }, (_, i) => ({
    email: `test.user.${i}@example.com`,
    firstName: `TestUser${i}`,
    lastName: `LastName${i}`,
  }));
});

export const options = {
  stages: [
    { duration: '30s', target: VUS },
    { duration: DURATION, target: VUS },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'],
    http_req_failed: ['rate<0.01'],
  },
};

function getRandomItem<T>(array: readonly T[]): T {
  return array[Math.floor(Math.random() * array.length)];
}

function withRetry<T>(fn: () => T, retries = MAX_RETRIES): T {
  for (let i = 0; i < retries; i++) {
    try {
      return fn();
    } catch (e) {
      if (i === retries - 1) throw e;
      console.log(`Attempt ${i + 1} failed, retrying in ${RETRY_INTERVAL}s...`);
      sleep(RETRY_INTERVAL);
    }
  }
  throw new Error('All retries failed');
}

function authenticateUser(user: TestUser): string {
  return withRetry(() => {
    const loginRes = http.post(`${BASE_URL}/auth/test`, JSON.stringify(user), {
      headers: { 'Content-Type': 'application/json' },
    });

    if (loginRes.status >= 400) {
      console.error(`Authentication failed: Status ${loginRes.status}, Body: ${loginRes.body}`);
    }

    const cookies = loginRes.headers['Set-Cookie'];
    if (!cookies) {
      console.error('No cookies received from authentication');
      return '';
    }

    const sessionCookie = cookies.toString().match(/storeit_session=([^;]+)/);
    if (!sessionCookie) {
      console.error('Session cookie not found in response');
      return '';
    }

    return sessionCookie[1];
  });
}

function getAuthHeaders(sessionToken: string, orgId?: string) {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    Cookie: `storeit_session=${sessionToken}`,
  };
  
  if (orgId) {
    headers['x-organization-id'] = orgId;
  }
  
  return headers;
}

function makeRequest(method: string, url: string, body: any = null, headers: any = {}) {
  return withRetry(() => {
    const options = { headers };
    const response = method === 'GET'
      ? http.get(url, options)
      : http.post(url, JSON.stringify(body), options);

    if (response.status >= 400) {
      console.error(`Request failed: ${method} ${url} - Status: ${response.status}, Body: ${response.body}`);
    }
    return response;
  });
}

function createItem(headers: any) {
  const itemName = `Item${randomString(6)}`;
  const itemRes = makeRequest('POST', `${BASE_URL}/items`, {
    name: itemName,
    description: `Test item description ${randomString(20)}`,
  }, headers);
  return JSON.parse(itemRes.body).data;
}

function createVariant(itemId: string, headers: any) {
  const variantName = `Variant${randomString(6)}`;
  const article = randomString(12, '0123456789');
  const variantRes = makeRequest('POST', `${BASE_URL}/items/${itemId}/variants`, {
    name: variantName,
    article: article,
  }, headers);
  return JSON.parse(variantRes.body).data;
}

function createCellGroup(unitId: string, headers: any) {
  const cellGroupName = `CellGroup${randomString(6)}`;
  const cellGroupRes = makeRequest('POST', `${BASE_URL}/cells-groups`, {
    name: cellGroupName,
    alias: `CG${randomString(3)}`,
    unitId: unitId,
  }, headers);
  return JSON.parse(cellGroupRes.body).data;
}

function createCell(cellGroupId: string, headers: any) {
  const cellRes = makeRequest('POST', `${BASE_URL}/cells-groups/${cellGroupId}/cells`, {
    alias: `C${randomString(3)}`,
    row: randomIntBetween(1, 5),
    level: randomIntBetween(1, 3),
    position: randomIntBetween(1, 10),
  }, headers);
  return JSON.parse(cellRes.body).data;
}

function createInstance(itemId: string, variantId: string, cellId: string, headers: any) {
  const instanceRes = makeRequest('POST', `${BASE_URL}/items/${itemId}/instances`, {
    variantId: variantId,
    cellId: cellId,
  }, headers);
  return JSON.parse(instanceRes.body).data;
}

function createTask(unitId: string, instanceId: string, sourceCellId: string, targetCellId: string, headers: any) {
  const taskRes = makeRequest('POST', `${BASE_URL}/tasks`, {
    name: `Task${randomString(6)}`,
    description: `Test task description ${randomString(20)}`,
    type: 'movement',
    unitId: unitId,
    items: [{
      instanceId: instanceId,
      targetCellId: targetCellId
    }]
  }, headers);

  if (taskRes.status !== 200) {
    throw new Error(`Failed to create task: ${taskRes.status} ${taskRes.body}`);
  }

  return JSON.parse(taskRes.body).data;
}

export default function () {
  try {
    const user = getRandomItem(users);
    const sessionToken = authenticateUser(user);
    const baseHeaders = getAuthHeaders(sessionToken);

    const orgName = `TestOrg${randomString(8)}`;
    const orgRes = makeRequest('POST', `${BASE_URL}/orgs`, {
      name: orgName,
      subdomain: orgName.toLowerCase(),
    }, baseHeaders);

    if (orgRes.status === 200) {
      const org = JSON.parse(orgRes.body).data;
      const headers = getAuthHeaders(sessionToken, org.id);

      const unitRes = makeRequest('POST', `${BASE_URL}/units`, {
        name: `Warehouse${randomString(6)}`,
        alias: `WH${randomString(3)}`,
        address: `Test Address ${randomString(10)}`,
      }, headers);

      if (unitRes.status === 200) {
        const unit = JSON.parse(unitRes.body).data;

        const storageGroup = makeRequest('POST', `${BASE_URL}/storage-groups`, {
          name: `Storage${randomString(6)}`,
          alias: `ST${randomString(3)}`,
          unitId: unit.id,
        }, headers);

        const cellGroup = createCellGroup(unit.id, headers);
        const sourceCell = createCell(cellGroup.id, headers);
        const targetCell = createCell(cellGroup.id, headers);

        const item = createItem(headers);
        const variant = createVariant(item.id, headers);

        // const instance = createInstance(item.id, variant.id, sourceCell.id, headers);

        // const task = createTask(unit.id, instance.id, sourceCell.id, targetCell.id, headers);

        makeRequest('GET', `${BASE_URL}/orgs`, null, headers);
        makeRequest('GET', `${BASE_URL}/units`, null, headers);
        makeRequest('GET', `${BASE_URL}/storage-groups`, null, headers);
        makeRequest('GET', `${BASE_URL}/items`, null, headers);
        makeRequest('GET', `${BASE_URL}/tasks`, null, headers);
        makeRequest('GET', `${BASE_URL}/cells-groups`, null, headers);
        makeRequest('GET', `${BASE_URL}/items/${item.id}/variants`, null, headers);
        // makeRequest('GET', `${BASE_URL}/items/${item.id}/instances`, null, headers);

        // makeRequest('PUT', `${BASE_URL}/items/${item.id}`, {
        //   name: `UpdatedItem${randomString(6)}`,
        //   description: `Updated description ${randomString(20)}`,
        // }, headers);

        // makeRequest('PUT', `${BASE_URL}/items/${item.id}/variants/${variant.id}`, {
        //   name: `UpdatedVariant${randomString(6)}`,
        //   article: randomString(12, '0123456789'),
        // }, headers);

        // makeRequest('PUT', `${BASE_URL}/cells-groups/${cellGroup.id}/cells/${sourceCell.id}`, {
        //   alias: `UC${randomString(3)}`,
        //   row: sourceCell.row,
        //   level: sourceCell.level,
        //   position: sourceCell.position,
        // }, headers);
      }
    }
  } catch (e) {
    console.error(`Test iteration failed: ${e.message}`);
  }

  sleep(randomIntBetween(1, 5));
}

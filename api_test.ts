import { expect, test } from "bun:test";
import { createClient, createConfig } from './src/client/client'
import { getV1StFirstRun } from './src/client'
let admin_api_token: string | undefined, user_api_token: string | undefined;
const client = createClient(createConfig({ baseUrl: 'http://localhost:3000' }));

test('Check if first run', async () => {
    const response = await getV1StFirstRun({ client });
    expect(response).toBeDefined();
    expect(response).toHaveProperty('data');
    expect(response.data).toHaveProperty('first_run');
    expect(response.data?.first_run).toBe(true);
})

test('Bootstrap admin token', async () => {
    const { postV1StBootstrap } = await import('./src/client');
    const response = await postV1StBootstrap({ client });
    expect(response).toBeDefined();
    expect(response).toHaveProperty('data');
    expect(response.data).toHaveProperty('access_key');
    admin_api_token = response.data?.access_key;
    expect(admin_api_token).toBeDefined();
    console.log(admin_api_token)
})

test('Check if first run after bootstrap', async () => {
    const response = await getV1StFirstRun({ client });
    expect(response).toBeDefined();
    expect(response).toHaveProperty('data');
    expect(response.data).toHaveProperty('first_run');
    expect(response.data?.first_run).toBe(false);
})

test('Create service token', async () => {
    const { postV1St } = await import('./src/client');
    const response = await postV1St({
        client,
        body: {
            name: 'user-token',
        },
        headers: {
            "x-api-token": admin_api_token as string
        }
    });
    expect(response).toBeDefined();
    expect(response).toHaveProperty('data');
    expect(response.data).toHaveProperty('access_key');
    user_api_token = response.data?.access_key;
    expect(user_api_token).toBeDefined();
    console.log(user_api_token, 'user');
});

test('Cannot create second admin token', async () => {
    const { postV1StBootstrap } = await import('./src/client');
    const response = await postV1StBootstrap({ client });
    expect(response).toBeDefined();
    expect(response).toHaveProperty('error');
});

test('Cannot create service token with an already existing name', async () => {
    const { postV1St } = await import('./src/client');
    const response = await postV1St({
        client,
        body: {
            name: 'user-token',
        },
        headers: {
            "x-api-token": admin_api_token as string
        }
    });
    expect(response).toBeDefined();
    expect(response).toHaveProperty('error');
});

test('Delete service token', async () => {
    const { deleteV1StById, postV1St } = await import('./src/client');
    const createResponse = await postV1St({
        client,
        body: {
            name: 'temp-token',
        },
        headers: {
            "x-api-token": admin_api_token as string
        }
    });
    expect(createResponse).toBeDefined();
    expect(createResponse).toHaveProperty('data');
    expect(createResponse.data).toHaveProperty('id');
    const token_id = createResponse.data?.id;
    expect(token_id).toBeDefined();
    const deleteResponse = await deleteV1StById({
        client,
        path: { id: token_id as string },
        headers: {
            "x-api-token": admin_api_token as string
        }
    });
    expect(deleteResponse).toBeDefined();
    expect(deleteResponse).toHaveProperty('data');
});



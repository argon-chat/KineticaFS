import { expect, test, describe } from "bun:test";
import { createClient, createConfig } from './src/client/client'
import { getV1StFirstRun } from './src/client'
let admin_api_token: string | undefined, user_api_token: string | undefined;
const client = createClient(createConfig({ baseUrl: 'http://localhost:3000' }));

describe('Authorization Tests', () => {
    test('Bucket endpoints return 401 without token', async () => {
        const { getV1Bucket, postV1Bucket, getV1BucketById, patchV1BucketById, deleteV1BucketById } = await import('./src/client');
        const listResponse = await getV1Bucket({
            client,
            headers: {} as any
        });
        expect(listResponse).toHaveProperty('error');
        expect(listResponse.response?.status).toBe(401);
        const createResponse = await postV1Bucket({
            client,
            headers: {} as any,
            body: {
                name: 'test-bucket',
                region: 'us-east-1',
                endpoint: 'https://s3.amazonaws.com',
                access_key: 'test-key',
                secret_key: 'test-secret',
                use_ssl: true,
                s3_provider: 'aws',
                storage_type: 0
            }
        });
        expect(createResponse).toHaveProperty('error');
        expect(createResponse.response?.status).toBe(401);
        const getResponse = await getV1BucketById({
            client,
            path: { id: 'fake-id' },
            headers: {} as any
        });
        expect(getResponse).toHaveProperty('error');
        expect(getResponse.response?.status).toBe(401);
        const updateResponse = await patchV1BucketById({
            client,
            path: { id: 'fake-id' },
            headers: {} as any,
            body: {
                name: 'updated-bucket',
                region: 'us-west-2',
                endpoint: 'https://s3.amazonaws.com',
                access_key: 'test-key',
                secret_key: 'test-secret',
                use_ssl: true,
                s3_provider: 'aws',
                storage_type: 0
            }
        });
        expect(updateResponse).toHaveProperty('error');
        expect(updateResponse.response?.status).toBe(401);
        const deleteResponse = await deleteV1BucketById({
            client,
            path: { id: 'fake-id' },
            headers: {} as any
        });
        expect(deleteResponse).toHaveProperty('error');
        expect(deleteResponse.response?.status).toBe(401);
    });

    test('Bucket endpoints return 401 with invalid token', async () => {
        const { getV1Bucket, postV1Bucket, getV1BucketById, patchV1BucketById, deleteV1BucketById } = await import('./src/client');
        const invalidToken = 'invalid-token-12345';
        const listResponse = await getV1Bucket({
            client,
            headers: { 'x-api-token': invalidToken }
        });
        expect(listResponse).toHaveProperty('error');
        expect(listResponse.response?.status).toBe(401);
        const createResponse = await postV1Bucket({
            client,
            headers: { 'x-api-token': invalidToken },
            body: {
                name: 'test-bucket',
                region: 'us-east-1',
                endpoint: 'https://s3.amazonaws.com',
                access_key: 'test-key',
                secret_key: 'test-secret',
                use_ssl: true,
                s3_provider: 'aws',
                storage_type: 0
            }
        });
        expect(createResponse).toHaveProperty('error');
        expect(createResponse.response?.status).toBe(401);
        const getResponse = await getV1BucketById({
            client,
            path: { id: 'fake-id' },
            headers: { 'x-api-token': invalidToken }
        });
        expect(getResponse).toHaveProperty('error');
        expect(getResponse.response?.status).toBe(401);
        const updateResponse = await patchV1BucketById({
            client,
            path: { id: 'fake-id' },
            headers: { 'x-api-token': invalidToken },
            body: {
                name: 'updated-bucket',
                region: 'us-west-2',
                endpoint: 'https://s3.amazonaws.com',
                access_key: 'test-key',
                secret_key: 'test-secret',
                use_ssl: true,
                s3_provider: 'aws',
                storage_type: 0
            }
        });
        expect(updateResponse).toHaveProperty('error');
        expect(updateResponse.response?.status).toBe(401);
        const deleteResponse = await deleteV1BucketById({
            client,
            path: { id: 'fake-id' },
            headers: { 'x-api-token': invalidToken }
        });
        expect(deleteResponse).toHaveProperty('error');
        expect(deleteResponse.response?.status).toBe(401);
    });

    test('File endpoints return 401 without token', async () => {
        const { postV1File, postV1FileByIdFinalize, deleteV1FileById } = await import('./src/client');
        const initiateResponse = await postV1File({
            client,
            headers: {} as any,
            body: {
                regionId: 'us-east-1',
                bucketCode: 'test-bucket'
            }
        });
        expect(initiateResponse).toHaveProperty('error');
        expect(initiateResponse.response?.status).toBe(401);
        const finalizeResponse = await postV1FileByIdFinalize({
            client,
            path: { id: 'fake-id' },
            headers: {} as any
        });
        expect(finalizeResponse).toHaveProperty('error');
        expect(finalizeResponse.response?.status).toBe(401);
        const deleteResponse = await deleteV1FileById({
            client,
            path: { id: 'fake-id' },
            headers: {} as any
        });
        expect(deleteResponse).toHaveProperty('error');
        expect(deleteResponse.response?.status).toBe(401);
    });

    test('File endpoints return 401 with invalid token', async () => {
        const { postV1File, postV1FileByIdFinalize, deleteV1FileById } = await import('./src/client');
        const invalidToken = 'invalid-token-12345';
        const initiateResponse = await postV1File({
            client,
            headers: { 'x-api-token': invalidToken },
            body: {
                regionId: 'us-east-1',
                bucketCode: 'test-bucket'
            }
        });
        expect(initiateResponse).toHaveProperty('error');
        expect(initiateResponse.response?.status).toBe(401);
        const finalizeResponse = await postV1FileByIdFinalize({
            client,
            path: { id: 'fake-id' },
            headers: { 'x-api-token': invalidToken }
        });
        expect(finalizeResponse).toHaveProperty('error');
        expect(finalizeResponse.response?.status).toBe(401);
        const deleteResponse = await deleteV1FileById({
            client,
            path: { id: 'fake-id' },
            headers: { 'x-api-token': invalidToken }
        });
        expect(deleteResponse).toHaveProperty('error');
        expect(deleteResponse.response?.status).toBe(401);
    });

    test('Upload endpoints return 401 without token', async () => {
        const { patchV1UploadByBlob } = await import('./src/client');
        const uploadResponse = await patchV1UploadByBlob({
            client,
            path: { blob: 'fake-blob-id' },
            headers: {} as any,
            body: [1, 2, 3, 4]
        });
        expect(uploadResponse).toHaveProperty('error');
        expect(uploadResponse.response?.status).toBe(401);
    });

    test('Upload endpoints return 401 with invalid token', async () => {
        const { patchV1UploadByBlob } = await import('./src/client');
        const invalidToken = 'invalid-token-12345';
        const uploadResponse = await patchV1UploadByBlob({
            client,
            path: { blob: 'fake-blob-id' },
            headers: { 'x-api-token': invalidToken },
            body: [1, 2, 3, 4]
        });
        expect(uploadResponse).toHaveProperty('error');
        expect(uploadResponse.response?.status).toBe(401);
    });

    test('Service token endpoints return 401 without token', async () => {
        const { getV1St, postV1St, getV1StById, deleteV1StById } = await import('./src/client');
        const listResponse = await getV1St({
            client,
            headers: {} as any
        });
        expect(listResponse).toHaveProperty('error');
        expect(listResponse.response?.status).toBe(401);
        const createResponse = await postV1St({
            client,
            headers: {} as any,
            body: { name: 'test-token' }
        });
        expect(createResponse).toHaveProperty('error');
        expect(createResponse.response?.status).toBe(401);
        const getResponse = await getV1StById({
            client,
            path: { id: 'fake-id' },
            headers: {} as any
        });
        expect(getResponse).toHaveProperty('error');
        expect(getResponse.response?.status).toBe(401);
        const deleteResponse = await deleteV1StById({
            client,
            path: { id: 'fake-id' },
            headers: {} as any
        });
        expect(deleteResponse).toHaveProperty('error');
        expect(deleteResponse.response?.status).toBe(401);
    });

    test('Service token endpoints return 401 with invalid token', async () => {
        const { getV1St, postV1St, getV1StById, deleteV1StById } = await import('./src/client');
        const invalidToken = 'invalid-token-12345';
        const listResponse = await getV1St({
            client,
            headers: { 'x-api-token': invalidToken }
        });
        expect(listResponse).toHaveProperty('error');
        expect(listResponse.response?.status).toBe(401);
        const createResponse = await postV1St({
            client,
            headers: { 'x-api-token': invalidToken },
            body: { name: 'test-token' }
        });
        expect(createResponse).toHaveProperty('error');
        expect(createResponse.response?.status).toBe(401);
        const getResponse = await getV1StById({
            client,
            path: { id: 'fake-id' },
            headers: { 'x-api-token': invalidToken }
        });
        expect(getResponse).toHaveProperty('error');
        expect(getResponse.response?.status).toBe(401);
        const deleteResponse = await deleteV1StById({
            client,
            path: { id: 'fake-id' },
            headers: { 'x-api-token': invalidToken }
        });
        expect(deleteResponse).toHaveProperty('error');
        expect(deleteResponse.response?.status).toBe(401);
    });

    test('Public endpoints work without token', async () => {
        const { getV1StFirstRun, postV1StBootstrap } = await import('./src/client');
        const firstRunResponse = await getV1StFirstRun({ client });
        expect(firstRunResponse).toBeDefined();
        expect(firstRunResponse).toHaveProperty('data');
        expect(firstRunResponse.data).toHaveProperty('first_run');
        expect(typeof firstRunResponse.data?.first_run).toBe('boolean');
        const bootstrapResponse = await postV1StBootstrap({ client });
        expect(bootstrapResponse).toBeDefined();
        if (bootstrapResponse.error) {
            expect(bootstrapResponse.response?.status).not.toBe(401);
        }
        expect(bootstrapResponse).toHaveProperty('data');
        expect(bootstrapResponse.data).toHaveProperty('access_key');
        admin_api_token = bootstrapResponse.data?.access_key;
        expect(admin_api_token).toBeDefined();
        console.log(admin_api_token)
    });

    test('Admin endpoints return 403 with non-admin token', async () => {
        const { getV1Bucket, postV1Bucket, getV1St, postV1St } = await import('./src/client');
        if (user_api_token) {
            const bucketResponse = await getV1Bucket({
                client,
                headers: { 'x-api-token': user_api_token }
            });
            expect(bucketResponse).toHaveProperty('error');
            expect([401, 403]).toContain(bucketResponse.response?.status);
            const tokenResponse = await getV1St({
                client,
                headers: { 'x-api-token': user_api_token }
            });
            expect(tokenResponse).toHaveProperty('error');
            expect([401, 403]).toContain(tokenResponse.response?.status);
        }
    });
});

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

test('Creates and deletes bucket successfully', async () => {
    const { postV1Bucket, deleteV1BucketById } = await import('./src/client');
    const createResponse = await postV1Bucket({
        client,
        headers: {
            "x-api-token": admin_api_token as string
        },
        body: {
            name: 'test-bucket',
            region: 'us-east-1',
            endpoint: 'https://s3.amazonaws.com',
            access_key: 'test-key',
            secret_key: 'test-secret',
            use_ssl: true,
            s3_provider: 'aws',
            storage_type: 0
        }
    });
    expect(createResponse).toBeDefined();
    expect(createResponse).toHaveProperty('data');
    expect(createResponse.data).toHaveProperty('id');
    const bucket_id = createResponse.data?.id;
    expect(bucket_id).toBeDefined();
    const deleteResponse = await deleteV1BucketById({
        client,
        path: { id: bucket_id as string },
        headers: {
            "x-api-token": admin_api_token as string
        }
    });
    expect(deleteResponse).toBeDefined();
    expect(deleteResponse).toHaveProperty('data');
});

test('Lists buckets after creating a bucket', async () => {
    const { postV1Bucket, getV1Bucket } = await import('./src/client');
    const createResponse = await postV1Bucket({
        client,
        headers: {
            "x-api-token": admin_api_token as string
        },
        body: {
            name: 'test-bucket',
            region: 'us-east-1',
            endpoint: 'https://s3.amazonaws.com',
            access_key: 'test-key',
            secret_key: 'test-secret',
            use_ssl: true,
            s3_provider: 'aws',
            storage_type: 0
        }
    });
    expect(createResponse).toBeDefined();
    expect(createResponse).toHaveProperty('data');
    expect(createResponse.data).toHaveProperty('id');
    const bucket_id = createResponse.data?.id;
    expect(bucket_id).toBeDefined();
    const listResponse = await getV1Bucket({
        client,
        headers: {
            "x-api-token": admin_api_token as string
        }
    });
    expect(listResponse).toBeDefined();
    expect(listResponse).toHaveProperty('data');
    const buckets = listResponse.data;
    expect(Array.isArray(buckets)).toBe(true);
    const createdBucket = buckets?.find((bucket: any) => bucket.id === bucket_id);
    expect(createdBucket).toBeDefined();
});

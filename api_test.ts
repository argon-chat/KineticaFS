import { expect, test, describe } from "bun:test";
import { createClient, createConfig } from './src/client/client'
import { firstRunCheck, ModelsBucket } from './src/client'
let admin_api_token: string | undefined, user_api_token: string | undefined;
const client = createClient(createConfig({ baseUrl: 'http://localhost:3000' }));

describe('Authorization Tests', () => {
    test('Bucket endpoints return 401 without token', async () => {
        const { listBuckets, createBucket, getBucket, updateBucket, deleteBucket } = await import('./src/client');
        const listResponse = await listBuckets({
            client,
            headers: {} as any
        });
        expect(listResponse).toHaveProperty('error');
        expect(listResponse.response?.status).toBe(401);
        const createResponse = await createBucket({
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
        const getResponse = await getBucket({
            client,
            path: { id: 'fake-id' },
            headers: {} as any
        });
        expect(getResponse).toHaveProperty('error');
        expect(getResponse.response?.status).toBe(401);
        const updateResponse = await updateBucket({
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
        const deleteResponse = await deleteBucket({
            client,
            path: { id: 'fake-id' },
            headers: {} as any
        });
        expect(deleteResponse).toHaveProperty('error');
        expect(deleteResponse.response?.status).toBe(401);
    });

    test('Bucket endpoints return 401 with invalid token', async () => {
        const { listBuckets, createBucket, getBucket, updateBucket, deleteBucket } = await import('./src/client');
        const invalidToken = 'invalid-token-12345';
        const listResponse = await listBuckets({
            client,
            headers: { 'x-api-token': invalidToken }
        });
        expect(listResponse).toHaveProperty('error');
        expect(listResponse.response?.status).toBe(401);
        const createResponse = await createBucket({
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
        const getResponse = await getBucket({
            client,
            path: { id: 'fake-id' },
            headers: { 'x-api-token': invalidToken }
        });
        expect(getResponse).toHaveProperty('error');
        expect(getResponse.response?.status).toBe(401);
        const updateResponse = await updateBucket({
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
        const deleteResponse = await deleteBucket({
            client,
            path: { id: 'fake-id' },
            headers: { 'x-api-token': invalidToken }
        });
        expect(deleteResponse).toHaveProperty('error');
        expect(deleteResponse.response?.status).toBe(401);
    });

    test('File endpoints return 401 without token', async () => {
        const { initiateFileUpload, finalizeFileUpload, deleteFile } = await import('./src/client');
        const initiateResponse = await initiateFileUpload({
            client,
            headers: {} as any,
            body: {
                regionId: 'us-east-1',
                bucketCode: 'test-bucket'
            }
        });
        expect(initiateResponse).toHaveProperty('error');
        expect(initiateResponse.response?.status).toBe(401);
        const finalizeResponse = await finalizeFileUpload({
            client,
            path: { blob: 'fake-id' },
            headers: {} as any
        });
        expect(finalizeResponse).toHaveProperty('error');
        expect(finalizeResponse.response?.status).toBe(401);
        const deleteResponse = await deleteFile({
            client,
            path: { id: 'fake-id' },
            headers: {} as any
        });
        expect(deleteResponse).toHaveProperty('error');
        expect(deleteResponse.response?.status).toBe(401);
    });

    test('File endpoints return 401 with invalid token', async () => {
        const { initiateFileUpload, finalizeFileUpload, deleteFile } = await import('./src/client');
        const invalidToken = 'invalid-token-12345';
        const initiateResponse = await initiateFileUpload({
            client,
            headers: { 'x-api-token': invalidToken },
            body: {
                regionId: 'us-east-1',
                bucketCode: 'test-bucket'
            }
        });
        expect(initiateResponse).toHaveProperty('error');
        expect(initiateResponse.response?.status).toBe(401);
        const finalizeResponse = await finalizeFileUpload({
            client,
            path: { blob: 'fake-id' },
            headers: { 'x-api-token': invalidToken }
        });
        expect(finalizeResponse).toHaveProperty('error');
        expect(finalizeResponse.response?.status).toBe(401);
        const deleteResponse = await deleteFile({
            client,
            path: { id: 'fake-id' },
            headers: { 'x-api-token': invalidToken }
        });
        expect(deleteResponse).toHaveProperty('error');
        expect(deleteResponse.response?.status).toBe(401);
    });

    test('Upload endpoints return 401 without token', async () => {
        const { uploadFileBlob } = await import('./src/client');
        const uploadResponse = await uploadFileBlob({
            client,
            path: { blob: 'fake-blob-id' },
            headers: {} as any,
            body: [1, 2, 3, 4]
        });
        expect(uploadResponse).toHaveProperty('error');
        expect(uploadResponse.response?.status).toBe(401);
    });

    test('Upload endpoints return 401 with invalid token', async () => {
        const { uploadFileBlob } = await import('./src/client');
        const invalidToken = 'invalid-token-12345';
        const uploadResponse = await uploadFileBlob({
            client,
            path: { blob: 'fake-blob-id' },
            headers: { 'x-api-token': invalidToken },
            body: [1, 2, 3, 4]
        });
        expect(uploadResponse).toHaveProperty('error');
        expect(uploadResponse.response?.status).toBe(401);
    });

    test('Service token endpoints return 401 without token', async () => {
        const { listAllServiceTokens, createServiceToken, getServiceToken, deleteServiceToken } = await import('./src/client');
        const listResponse = await listAllServiceTokens({
            client,
            headers: {} as any
        });
        expect(listResponse).toHaveProperty('error');
        expect(listResponse.response?.status).toBe(401);
        const createResponse = await createServiceToken({
            client,
            headers: {} as any,
            body: { name: 'test-token' }
        });
        expect(createResponse).toHaveProperty('error');
        expect(createResponse.response?.status).toBe(401);
        const getResponse = await getServiceToken({
            client,
            path: { id: 'fake-id' },
            headers: {} as any
        });
        expect(getResponse).toHaveProperty('error');
        expect(getResponse.response?.status).toBe(401);
        const deleteResponse = await deleteServiceToken({
            client,
            path: { id: 'fake-id' },
            headers: {} as any
        });
        expect(deleteResponse).toHaveProperty('error');
        expect(deleteResponse.response?.status).toBe(401);
    });

    test('Service token endpoints return 401 with invalid token', async () => {
        const { listAllServiceTokens, createServiceToken, getServiceToken, deleteServiceToken } = await import('./src/client');
        const invalidToken = 'invalid-token-12345';
        const listResponse = await listAllServiceTokens({
            client,
            headers: { 'x-api-token': invalidToken }
        });
        expect(listResponse).toHaveProperty('error');
        expect(listResponse.response?.status).toBe(401);
        const createResponse = await createServiceToken({
            client,
            headers: { 'x-api-token': invalidToken },
            body: { name: 'test-token' }
        });
        expect(createResponse).toHaveProperty('error');
        expect(createResponse.response?.status).toBe(401);
        const getResponse = await getServiceToken({
            client,
            path: { id: 'fake-id' },
            headers: { 'x-api-token': invalidToken }
        });
        expect(getResponse).toHaveProperty('error');
        expect(getResponse.response?.status).toBe(401);
        const deleteResponse = await deleteServiceToken({
            client,
            path: { id: 'fake-id' },
            headers: { 'x-api-token': invalidToken }
        });
        expect(deleteResponse).toHaveProperty('error');
        expect(deleteResponse.response?.status).toBe(401);
    });

    test('Public endpoints work without token', async () => {
        const { firstRunCheck, bootstrapAdminToken } = await import('./src/client');
        const firstRunResponse = await firstRunCheck({ client });
        expect(firstRunResponse).toBeDefined();
        expect(firstRunResponse).toHaveProperty('data');
        expect(firstRunResponse.data).toHaveProperty('first_run');
        expect(typeof firstRunResponse.data?.first_run).toBe('boolean');
        const bootstrapResponse = await bootstrapAdminToken({ client });
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
        const { listBuckets, createBucket, listAllServiceTokens, createServiceToken } = await import('./src/client');
        if (user_api_token) {
            const bucketResponse = await listBuckets({
                client,
                headers: { 'x-api-token': user_api_token }
            });
            expect(bucketResponse).toHaveProperty('error');
            expect([401, 403]).toContain(bucketResponse.response?.status);
            const tokenResponse = await listAllServiceTokens({
                client,
                headers: { 'x-api-token': user_api_token }
            });
            expect(tokenResponse).toHaveProperty('error');
            expect([401, 403]).toContain(tokenResponse.response?.status);
        }
    });
});

test('Check if first run after bootstrap', async () => {
    const response = await firstRunCheck({ client });
    expect(response).toBeDefined();
    expect(response).toHaveProperty('data');
    expect(response.data).toHaveProperty('first_run');
    expect(response.data?.first_run).toBe(false);
})

test('Create service token', async () => {
    const { createServiceToken } = await import('./src/client');
    const response = await createServiceToken({
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
    const { bootstrapAdminToken } = await import('./src/client');
    const response = await bootstrapAdminToken({ client });
    expect(response).toBeDefined();
    expect(response).toHaveProperty('error');
});

test('Cannot create service token with an already existing name', async () => {
    const { createServiceToken } = await import('./src/client');
    const response = await createServiceToken({
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
    const { deleteServiceToken, createServiceToken } = await import('./src/client');
    const createResponse = await createServiceToken({
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
    const deleteResponse = await deleteServiceToken({
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
    const { createBucket, deleteBucket } = await import('./src/client');
    const createResponse = await createBucket({
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
    const deleteResponse = await deleteBucket({
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
    const { createBucket, listBuckets, deleteBucket } = await import('./src/client');
    const createResponse = await createBucket({
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
    const listResponse = await listBuckets({
        client,
        headers: {
            "x-api-token": admin_api_token as string
        }
    });
    expect(listResponse).toBeDefined();
    expect(listResponse).toHaveProperty('data');
    const buckets = listResponse.data;
    expect(Array.isArray(buckets)).toBe(true);
    const createdBucket = buckets?.find((bucket: ModelsBucket) => bucket.id === bucket_id);
    expect(createdBucket).toBeDefined();
    const deleteResponse = await deleteBucket({
        client,
        path: { id: bucket_id as string },
        headers: {
            "x-api-token": admin_api_token as string
        }
    });
    expect(deleteResponse).toBeDefined();
    expect(deleteResponse).toHaveProperty('data');
});

test('creates the actual buckets from compose file', async () => {
    const { createBucket } = await import('./src/client');
    const region1Response = await createBucket({
        client,
        headers: {
            "x-api-token": admin_api_token as string
        },
        body: {
            name: 'region1-storage',
            region: 'region1',
            endpoint: 'http://localhost:8333',
            access_key: 'argon',
            secret_key: 'argon',
            use_ssl: false,
            s3_provider: 'seaweedfs',
            storage_type: 0
        }
    });

    const region2Response = await createBucket({
        client,
        headers: {
            "x-api-token": admin_api_token as string
        },
        body: {
            name: 'region2-storage',
            region: 'region2',
            endpoint: 'http://localhost:8334',
            access_key: 'argon',
            secret_key: 'argon',
            use_ssl: false,
            s3_provider: 'seaweedfs',
            storage_type: 0
        }
    });

    expect(region1Response).toBeDefined();
    expect(region1Response).toHaveProperty('data');
    expect(region2Response).toBeDefined();
    expect(region2Response).toHaveProperty('data');
    expect(region1Response.data).toHaveProperty('id');
    expect(region2Response.data).toHaveProperty('id');
    console.log('Region 1 bucket created:', region1Response.data?.id);
    console.log('Region 2 bucket created:', region2Response.data?.id);
})
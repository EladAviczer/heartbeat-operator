export const options = {
    insecureSkipTLSVerify: true,
    scenarios: {
        control_plane_stress: {
            executor: 'constant-vus',
            vus: 10,
            duration: '30s',
        },
    },
    thresholds: {
        http_req_failed: ['rate<0.05'], // less than 5% errors
        http_req_duration: ['p(95)<500'], // 95% of requests under 500ms
    },
};

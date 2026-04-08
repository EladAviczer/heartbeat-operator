import { check, sleep } from 'k6';
import { options } from '../config.js';
import { createProbe } from '../utils/k8s.js';

export { options };

const K8S_API_URL = __ENV.KUBERNETES_API_URL;
const K8S_TOKEN = __ENV.KUBERNETES_TOKEN;

export function setup() {
    if (!K8S_API_URL || !K8S_TOKEN) {
        throw new Error("KUBERNETES_API_URL and KUBERNETES_TOKEN environment variables must be defined");
    }
}

export default function () {
    // Add timestamp to make probes unique across multiple load test runs
    const probeName = `load-test-${Date.now()}-${__VU}-${__ITER}`;
    const res = createProbe(K8S_API_URL, K8S_TOKEN, probeName, "default");
    
    check(res, {
        'probe created successfully (201)': (r) => r.status === 201,
    });
    
    sleep(0.1); // Small delay to pace requests
}

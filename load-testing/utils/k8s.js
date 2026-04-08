import http from 'k6/http';

export function createProbe(apiHost, token, name, namespace) {
    const payload = JSON.stringify({
        apiVersion: "probes.ready.io/v1alpha1",
        kind: "Probe",
        metadata: {
            name: name,
            namespace: namespace,
            labels: {
                "created-by": "k6-load-test"
            }
        },
        spec: {
            checkType: "http",
            checkTarget: "http://dummy-target.default.svc.cluster.local",
            interval: "10s",
            timeout: "2s"
        }
    });

    const params = {
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
        },
    };

    return http.post(`${apiHost}/apis/probes.ready.io/v1alpha1/namespaces/${namespace}/probes`, payload, params);
}

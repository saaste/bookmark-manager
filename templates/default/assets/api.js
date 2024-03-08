export const getMetadata = (baseUrl, targetUrl) => {
    let url = `${baseUrl}/api/metadata?url=${targetUrl}`
    return fetch(url, { method: 'GET' })
    .then(resp => {
        if (!resp.ok) {
            return Promise.resolve({});
        }
        return resp.json();
    })
    .then(data => {
        return data;
    })
    .catch(err => { console.log("Fetching error", err) });
}

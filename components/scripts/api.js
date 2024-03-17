export const getMetadata = (baseUrl, targetUrl) => {
    let url = `${baseUrl}api/metadata?url=${targetUrl}`
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
    .catch(err => { console.log("Error fetching metadata", err) });
}

export const getTags = (baseUrl) => {
    let url = `${baseUrl}api/tags`
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
    .catch(err => { console.log("Error fetching tags", err) });
}
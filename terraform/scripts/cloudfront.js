function handler(event) {
    var request = event.request;
    var uri = request.uri;

    // Check if URI lacks a file extension and does not end in /
    if (uri.length > 1) {
        if (!uri.includes('.')) {
            if (uri.endsWith('/')) {
                request.uri += 'index.html';
            } else {
                request.uri += '.html';
            }
        }
    }
    return request;
}
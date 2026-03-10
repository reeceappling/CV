function handler(event) {
    var request = event.request;
    var uri = request.uri;

    // Check if URI lacks a file extension and does not end in /
    if (uri.length > 1 && !uri.endsWith('/') && !uri.includes('.')) {
        request.uri += '.html';
    }
    return request;
}
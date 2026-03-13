function handler(event) {
    var request = event.request;
    var uri = request.uri;

    // Check if URI lacks a file extension
    if (uri.length > 1) {
        if (!uri.includes('.')) {
            // If directory, go to folder index page, otherwise just go to the named page
            if (uri.endsWith('/')) {
                request.uri += 'index.html';
            } else {
                request.uri += '.html';
            }
        }
    }
    return request;
}
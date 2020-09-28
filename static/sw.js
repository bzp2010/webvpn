self.version = '2';

/*self.addEventListener('install', function(event) {
    event.waitUntil(Promise.resolve());
});

self.addEventListener('activate', function(e) {
    console.log('Activate event');
});*/

self.addEventListener('fetch', function (e) {
    /*let path = e.request.url.split('/');

    for (let i = 0; i < 3; i++) path.shift();
    path = "/" + path.join('/');

    if (path !== self.GlobalData.path) {
        e.respondWith(fetch("/" + self.GlobalData.service.name + path));
    } else {
        e.respondWith(fetch(path));
    }*/

    /*if (event.request.url === 'http://127.0.0.1:8080/data.txt') {
        event.respondWith(new Response('Hello World!'))
    }*/
});

self.addEventListener('message', function(event){
    self.GlobalData = {}
    self.GlobalData.path = event.data.path;
    self.GlobalData.service = {
        name: event.data.service.name
    }
});
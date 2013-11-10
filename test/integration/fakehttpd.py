import BaseHTTPServer
import os
import threading


COMMITS = open(os.path.join(os.path.dirname(__file__), '../data/commits.json')).read()


class Handler(BaseHTTPServer.BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == '/repos/mhrivnak/testing/pulls/3/commits':
            self.send_response(200)
            self.end_headers()
            self.wfile.write(COMMITS)
        else:
            self.send_response(404)


def run():
    server = BaseHTTPServer.HTTPServer(('localhost', 8081), Handler)
    threading.Thread(target=server.serve_forever).start()
    return server
            

from collections import namedtuple
import smtpd
import threading


Message = namedtuple('Message', ['peer', 'mailfrom', 'rcpttos', 'data'])


class FakeSMTPServer(smtpd.SMTPServer):
    def __init__(self):
        # this is an old-style class :(
        smtpd.SMTPServer.__init__(self, ('127.0.0.1', 1025), None)
        self._messages = []
        self.lock = threading.Lock()

    def process_message(self, peer, mailfrom, rcpttos, data):
        with self.lock:
            self.messages.append(Message(peer, mailfrom, rcpttos, data))

    def clear_messages(self):
        with self.lock:
            self.messages = []

    @property
    def messages(self):
        with self.lock:
            return self._messages.copy()

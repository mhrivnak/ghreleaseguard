import asyncore
import os
import smtpd
import subprocess
import threading
import time
import unittest
import urllib2
import urlparse

import fakesmtpd
import fakehttpd


DATADIR = os.path.join(os.path.dirname(__file__), '../data/')
CONFIGPATH = os.path.join(DATADIR, 'ghreleaseguard.conf')
GOPATH = os.getenv('GOPATH')
assert GOPATH is not None
BINFILE = os.path.join(GOPATH, 'bin/ghreleaseguard')
BASEURL = 'http://localhost:8080/api/v1/'


class TestGHRG(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        # start GHRG service
        cls.proc = cls._run_ghrg()

        # listen for mail
        cls.mailserver = fakesmtpd.FakeSMTPServer()
        cls.loopthread = threading.Thread(target=asyncore.loop, kwargs={'timeout':1})
        cls.loopthread.start()

        # start our fake github API
        cls.httpd = fakehttpd.run()

        # let services get started
        time.sleep(1)

    @classmethod
    def tearDownClass(cls):
        cls.mailserver.close()
        cls.httpd.shutdown()
        cls.proc.terminate()
        print cls.proc.stdout.read()
        print cls.proc.stderr.read()

    @staticmethod
    def _run_ghrg():
        # build it
        returncode = subprocess.call(['go', 'install', 'github.com/mhrivnak/ghreleaseguard'])
        assert returncode == 0
        
        # run it
        env = os.environ.copy()
        env['GHRGCONFIGPATH'] = CONFIGPATH
        return subprocess.Popen([BINFILE], env=env, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        
    def setUp(self):
        self.mailserver.clear_messages()

    def test_good_push(self):
        """
        Hits the API with JSON from a push that does not include the
        forbidden commit. This should not produce any notifications.
        """
        url = urlparse.urljoin(BASEURL, 'push')
        data = open(os.path.join(DATADIR, 'goodpush.json')).read()
        response = urllib2.urlopen(url, data)
        
        self._validate_good_response(response)

    def test_bad_push(self):
        """
        Hits the API with JSON from a push that does include the
        forbidden commit. This should produce an email.
        """
        url = urlparse.urljoin(BASEURL, 'push')
        data = open(os.path.join(DATADIR, 'badpush.json')).read()
        response = urllib2.urlopen(url, data)

        self._validate_bad_response(response)

    def test_good_pr(self):
        url = urlparse.urljoin(BASEURL, 'pullrequest')
        data = open(os.path.join(DATADIR, 'goodpr.json')).read()
        response = urllib2.urlopen(url, data)
        
        self._validate_good_response(response)

    def test_bad_pr(self):
        url = urlparse.urljoin(BASEURL, 'pullrequest')
        data = open(os.path.join(DATADIR, 'badpr.json')).read()
        response = urllib2.urlopen(url, data)

        self._validate_bad_response(response)

    def _validate_good_response(self, response):
        # wait for GHRG to process the request
        time.sleep(.2)

        self.assertEqual(response.getcode(), 200)
        self.assertEqual(len(self.mailserver.messages), 0)

    def _validate_bad_response(self, response):
        # wait for GHRG to process the request
        time.sleep(.2)

        self.assertEqual(response.getcode(), 200)
        self.assertEqual(len(self.mailserver.messages), 1)
        message = self.mailserver.messages[0]
        self.assertEqual(message.mailfrom, 'testsender@hrivnak.org')
        self.assertEqual(message.rcpttos, ['testreceiver@hrivnak.org'])

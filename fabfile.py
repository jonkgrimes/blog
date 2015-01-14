from __future__ import with_statement
from fabric.api import *
from fabric.contrib.console import confirm

env.hosts = ["deploy@jonkgrimes.com"]
env.forward_agent = True

code_dir = '/var/www/blog'

def deploy():
    with settings(warn_only=True):
        if run("test -d %s" % code_dir).failed:
            sudo("mkdir %s" % code_dir)
            sudo("chown deploy %s" % code_dir)
            run("git clone git@github.com:jonkgrimes/blog.git %s" % code_dir)
    with cd(code_dir):
        run("git pull")
        run("go build")
        run("./blog")

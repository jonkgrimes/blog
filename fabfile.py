from fabric.api import local

def prepare_deploy():
    local("go build")

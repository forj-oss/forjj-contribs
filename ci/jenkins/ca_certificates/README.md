// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
# Introduction

During your image build process, you can add any other root certificate, like HP Entreprise root certificate. This will ensure any https connection to github entreprise will be supported, out of the box.

You can also accept insecure connection, but not recommended for production.

To add your CA certificates, copy it in this directory and rebuild your image with bin/build.sh
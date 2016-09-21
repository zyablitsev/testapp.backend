#
# testapp.backend Dockerfile
#

# Pull base image.
FROM scratch

MAINTAINER Vasiliy Zyablitsev <vasiliy.zyablitsev@gmail.com>

# Set instructions on build.
ADD testapp /
ADD testappctl /
ADD stuff/config.yaml /etc/testapp/

# Expose ports.
EXPOSE 8701

# Run app
WORKDIR /
CMD ["/testapp"]

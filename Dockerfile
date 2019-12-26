FROM centos:7

ARG user

# The following packages are required for the following reasons:
#  - ccache (better for incremental builds)
#  - docbook-dtds for asciidoc's validation
#  - epel-release (for ccache)
#  - ghc-compiler (to bootstrap ghc)
#  - pexpect (rr)
#  - perl-ExtUtils-MakeMaker (git)
#  TODO: Remove cmake, we have it now (self-contained, too).
RUN yum --assumeyes update \
    && yum --assumeyes install \
         epel-release \
    && yum --assumeyes install \
         @'Development Tools' \
         glibc-static \
         libstdc++-static \
         ccache \
         cmake \
         docbook-dtds \
				 ghc-compiler \
         man-db \
         perl-ExtUtils-MakeMaker \
         pexpect \
         python \
         sudo \
    && yum clean all
# To test if the binaries are hardcoded to the build user's path
RUN useradd $user --create-home
# Allows us to install packages while testing.
RUN echo "$user ALL=(ALL) NOPASSWD: ALL" | EDITOR='tee -a' visudo
USER $user
WORKDIR /home/$user
ENV PATH /home/$user/.local/bin:/usr/lib64/ccache:$PATH

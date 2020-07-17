FROM scratch
ADD ./keepopendeleted /
ENTRYPOINT ["/keepopendeleted"]
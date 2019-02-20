eval 'set +o history' 2>/dev/null || setopt HIST_IGNORE_SPACE 2>/dev/null
 touch ~/.gitcookies
 chmod 0600 ~/.gitcookies

 git config --global http.cookiefile ~/.gitcookies

 tr , \\t <<\__END__ >>~/.gitcookies
go.googlesource.com,FALSE,/,TRUE,2147483647,o,git-hat3d111.gmail.com=1/mKVZHOD5A3D_plZdKA-SpSM11eQWkZaSBduD06L9Qf4
go-review.googlesource.com,FALSE,/,TRUE,2147483647,o,git-hat3d111.gmail.com=1/mKVZHOD5A3D_plZdKA-SpSM11eQWkZaSBduD06L9Qf4
__END__
eval 'set -o history' 2>/dev/null || unsetopt HIST_IGNORE_SPACE 2>/dev/null
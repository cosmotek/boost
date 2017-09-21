-- this file is an example lua module for use
-- with boost config files. it simply provides
-- a method for bumping semver strings.

-- regex module provided by boost
local re = require("re")
local semver = {}

semver.rc = "rc"
semver.alpha = "alpha"
semver.beta = "beta"

semver.major = 0
semver.minor = 1
semver.patch = 2

function semver.bump(original, bump_type)
  return original .. bump_type
end

return semver

-- this is a simple demo module for reading
-- /writing files with strings.
local fs = {}

function fs.readstr(filename)
  local f = io.open(filename, "rb")
  local content = f:read("*all")
  f:close()
  return content
end

-- r = read
-- w = write or create
-- a = append
-- r+ = read/write
-- w+ = overwrite with r/w permissions
-- a+ = r/w append or create

function fs.writestr(filename, string)
  local f = io.open(filename, "w+")
  io.output(f)
  io.write(string)
  io.close(f)
end

return fs

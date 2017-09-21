local str = {}
local charset = {}

-- qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890
for i = 48,  57 do table.insert(charset, string.char(i)) end
for i = 65,  90 do table.insert(charset, string.char(i)) end
for i = 97, 122 do table.insert(charset, string.char(i)) end

function str.random(length)
  math.randomseed(os.time())

  if length > 0 then
    return str.random(length - 1) .. charset[math.random(1, #charset)]
  else
    return ""
  end
end

return str

local md = require "markdown"

local fout, err, errno = io.open("./index.html", "w")
if err ~= nil then
  error("error opening index.html: ", errno, ": ", err)
end

local html, err = md.dofile("./README.md")
if err ~= nil then
  error(err)
end

fout:write(html)
fout:close()

local M = {}

function M.iso_datetime()
    return os.date("!%Y-%m-%dT%H:%M:%SZ")
end

return M

local M = {}

local list = require("fedilist.views.list")

function M.setup(opts)
    vim.api.nvim_create_user_command("GetList", function(opts)
        list.display_by_id(0)
    end, {})
end

return M

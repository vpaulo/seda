# External String Utils Module
module StringUtils ::
  fn concat(str1, str2) ::
    str1 + str2
  end

  fn repeat(str, times) ::
    var result = ""
    for i in [1, 2, 3, 4, 5] ::
        if i <= times ::
        result = result + str
        end
    end
    return result
  end

  fn addPrefix(str) ::
    "external_" + str
  end
end

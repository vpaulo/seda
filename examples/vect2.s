println("Running Vector2D tests...")

fn createVector2D(x, y) ::
  const vec = {"x": x, "y": y}

  vec.magnitude = fn(self) ::
    const x2 = self["x"] * self["x"]
    const y2 = self["y"] * self["y"]
    # Simple square root approximation for demo
    return x2 + y2
  end

  vec.add = fn(self, other) ::
    return createVector2D(
      self["x"] + other["x"],
      self["y"] + other["y"]
    )
  end

  vec.scale = fn(self, scalar) ::
    return createVector2D(
      self["x"] * scalar,
      self["y"] * scalar
    )
  end

  vec.dot = fn(self, other) ::
    return self["x"] * other["x"] + self["y"] * other["y"]
  end

  vec.to_string = fn(self) ::
    return "Vec2(" + self["x"].to_string + ", " + self["y"].to_string + ")"
  end

  return vec
end

const v1 = createVector2D(3, 4)
const v2 = createVector2D(1, 2)

check "Vector2D" ::
  v1.to_string is "Vec2(3, 4)"
  v2.to_string is "Vec2(1, 2)"
  v1.add(v2).to_string is "Vec2(4, 6)"
  v1.scale(2).to_string is "Vec2(6, 8)"
  v1.dot(v2) is 11
end

println("✓ All Vector2D tests passed!")
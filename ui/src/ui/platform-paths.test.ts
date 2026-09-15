import { describe, expect, it } from "vitest";
import { WIN_LIMNO_APPDATA, WIN_LIMNO_WORKSPACE } from "./platform-paths.ts";

describe("platform-paths", () => {
  it("preserves Windows backslashes for display", () => {
    expect(WIN_LIMNO_WORKSPACE).toBe("%APPDATA%\\linmo\\workspace");
    expect(WIN_LIMNO_APPDATA).toBe("%APPDATA%\\linmo");
  });
});

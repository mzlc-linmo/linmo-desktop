/** Windows 路径展示文案（String.raw 保留反斜杠，避免 \\ 与 \\o 转义问题）。 */
export const WIN_LIMNO_APPDATA = String.raw`%APPDATA%\linmo`;
export const WIN_LIMNO_WORKSPACE = String.raw`%APPDATA%\linmo\workspace`;

export const UNIX_LIMNO_WORKSPACE = "~/.linmo/workspace";

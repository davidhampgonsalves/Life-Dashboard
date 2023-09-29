import fontforge
F = fontforge.open("../Downloads/Noto_Emoji/static/NotoEmoji-Regular.ttf")
for name in F:
    filename = "noto-emoji/emoji_u" + hex(F[name].unicode)[2:] + ".png"
    2660
    F[name].export(filename, 32)

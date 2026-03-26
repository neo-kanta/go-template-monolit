import os
import glob
import re

count = 0
for fn in glob.glob('services/fnd/fndm*/mapper.go'):
    with open(fn, 'r', encoding='utf-8') as f:
        content = f.read()

    def repl(m):
        full = m.group(0)
        field = m.group(1).lower()
        
        # Determine decimal places based on Thai Localization conventions
        if any(x in field for x in ['unit', 'rate', 'nav', 'share', 'stay']):
            return f"{full}.RoundBank(4)"
        else:
            return f"{full}.RoundBank(2)"

    # Match decimal.NewFromFloat(m.GetXxx()) or s.GetXxx()
    new_content = re.sub(r'decimal\.NewFromFloat\([a-zA-Z0-9_]+\.Get([A-Za-z0-9_]+)\(\)\)', repl, content)

    if new_content != content:
        with open(fn, 'w', encoding='utf-8') as f:
            f.write(new_content)
        count += 1
        print(f"Updated {fn}")

print(f"Total files updated: {count}")

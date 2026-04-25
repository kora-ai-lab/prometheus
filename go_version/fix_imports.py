import os
import re
import glob

def fix_file(filepath):
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Fix import block format:
        # Current malformed: import ("context"\n\t"encoding/json"...)
        # Should be: import (\n\t"context"\n\t"encoding/json"\n)
        
        # Step 1: Add newline after import (
        new_content = re.sub(r'^import \(\t"', 'import (\n    "', content, flags=re.MULTILINE)
        
        # Step 2: Fix other tab-indented imports  
        new_content = re.sub(r'^\t"', '    "', new_content, flags=re.MULTILINE)
        
        # Step 3: Add closing ) and blank line before type/func
        # Find last import that isn't followed by another import or ) 
        # Pattern: last import line, then on next line type/func
        new_content = re.sub(r'^    "([a-zA-Z0-9_\-/.]+)"\n\n^(type|func)', r'    "\1"\n)\n\n\2', new_content, flags=re.MULTILINE)
        
        if new_content != content:
            with open(filepath, 'w', encoding='utf-8') as f:
                f.write(new_content)
            return True
        return False
    except Exception as e:
        print(f'Error in {filepath}: {e}')
        return False

count = 0
for filepath in glob.glob('C:/Users/junio/OneDrive/AI AGENT HACKATON/Prometheus/go_version/**/*.go', recursive=True):
    if fix_file(filepath):
        count += 1
        print(f'Fixed: {filepath}')

print(f'Total files re-fixed: {count}')
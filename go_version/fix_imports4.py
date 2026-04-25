import os
import re
import glob

def fix_file(filepath):
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Pattern: import followed by all imports without proper formatting
        # import ( + import list + no )
        
        result_lines = []
        
        # Find import block and rebuild properly
        match = re.search(r'^import \(\s*$', content, re.MULTILINE)
        if match:
            # Find where imports are
            start = content.find('import (')
            # Find where the imports end (first line that doesn't start with " or tab"
            rest = content[start:]
            import_lines = []
            
            lines = rest.split('\n')
            i = 0
            while i < len(lines):
                line = lines[i].strip()
                if line == '':
                    i += 1
                    continue
                if line.startswith('"'):
                    import_lines.append(line.rstrip(','))
                    i += 1
                else:
                    break
            
            # Now rebuild
            before = content[:start]
            result_lines.append(before)
            result_lines.append('import (')
            for imp in import_lines:
                result_lines.append('    ' + imp + ',')
            result_lines.append(')')
            result_lines.append('')
            if i < len(lines):
                result_lines.extend(lines[i:])
            
            new_content = '\n'.join(result_lines)
            
            if new_content != content:
                with open(filepath, 'w', encoding='utf-8') as f:
                    f.write(new_content)
                return True
        return False
    except Exception as e:
        print(f'Error in {filepath}: {e}')
        return False

count = 0
fixed_files = []
for filepath in glob.glob('C:/Users/junio/OneDrive/AI AGENT HACKATON/Prometheus/go_version/**/*.go', recursive=True):
    if fix_file(filepath):
        count += 1
        fixed_files.append(filepath)

print(f'Total files fixed: {count}')
print(f'First 10: {fixed_files[:10]}')
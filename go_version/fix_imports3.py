import os
import re
import glob

def fix_file(filepath):
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
        
        lines = content.split('\n')
        new_lines = []
        i = 0
        while i < len(lines):
            line = lines[i]
            
            # Detect import block start
            if line.strip() == 'import (':
                new_lines.append('import (')
                i += 1
                # Collect import lines
                while i < len(lines):
                    curr_line = lines[i].strip()
                    # If we've reached end of import block
                    if curr_line == '' or not curr_line.startswith('"'):
                        break
                    
                    # Process each import line
                    import_line = lines[i].strip()
                    if import_line.endswith(','):
                        new_lines.append('    ' + import_line)
                    else:
                        new_lines.append('    ' + import_line + ',')
                    i += 1
                
                # Close import block  
                if new_lines and new_lines[-1].strip().endswith(','):
                    new_lines.append(')')
                else:
                    new_lines[-1] = new_lines[-1].rstrip(',')
                    new_lines.append(')')
                
                new_lines.append('')  # blank line
                continue
            
            new_lines.append(line)
            i += 1
        
        new_content = '\n'.join(new_lines)
        
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
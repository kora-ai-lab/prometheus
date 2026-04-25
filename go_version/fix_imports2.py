import os
import re
import glob

def fix_file(filepath):
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
        
        lines = content.split('\n')
        new_lines = []
        in_import = False
        import_block_lines = []
        
        for i, line in enumerate(lines):
            if line.strip().startswith('import ('):
                in_import = True
                # Replace import ( with import (
                new_lines.append('import (')
                continue
            
            if in_import:
                # Check if this line ends the import block
                stripped = line.strip()
                
                # First, collect lines that look like imports
                if stripped.startswith('"'):
                    # Import line - add comma if not already there
                    if not stripped.endswith('"') and not stripped.endswith(','):
                        # Already properly quoted, just add to list
                        import_block_lines.append(line.strip())
                    elif stripped.endswith('"'):
                        # Add comma
                        import_block_lines.append(stripped + ',')
                    else:
                        import_block_lines.append(line.strip())
                    continue
                
                # Check for type/func which marks end of import block
                if (stripped.startswith('type ') or stripped.startswith('func ') or 
                    stripped.startswith('var ') or stripped.startswith('const ')):
                    # End of import block - write collected imports with proper formatting
                    for imp in import_block_lines:
                        new_lines.append('    ' + imp)
                    new_lines.append(')')
                    new_lines.append('')  # blank line
                    new_lines.append(line)
                    in_import = False
                    import_block_lines = []
                    continue
                
                # Otherwise, line doesn't look like import - end import block
                if import_block_lines:
                    for imp in import_block_lines:
                        new_lines.append('    ' + imp)
                    new_lines.append(')')
                    new_lines.append('')
                new_lines.append(line)
                in_import = False
                import_block_lines = []
                continue
            
            new_lines.append(line)
        
        # Handle files where import block wasn't detected
        if import_block_lines:
            for imp in import_block_lines:
                new_lines.append('    ' + imp)
            new_lines.append(')')
        
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
print('Files:', fixed_files[:10] if len(fixed_files) > 10 else fixed_files)
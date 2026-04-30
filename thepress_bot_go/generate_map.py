import os
import re

def parse_go_file(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    package = re.search(r'^package\s+(\w+)', content, re.MULTILINE)
    package_name = package.group(1) if package else "unknown"

    structs = re.findall(r'^type\s+(\w+)\s+struct', content, re.MULTILINE)
    interfaces = re.findall(r'^type\s+(\w+)\s+interface', content, re.MULTILINE)

    functions = []
    # match standard functions and methods
    func_pattern = re.compile(r'^func\s+(?:\(([^)]+)\)\s+)?(\w+)\s*\(', re.MULTILINE)
    for match in func_pattern.finditer(content):
        receiver = match.group(1)
        name = match.group(2)
        if receiver:
            functions.append(f"({receiver}) {name}")
        else:
            functions.append(name)

    return {
        "package": package_name,
        "structs": structs,
        "interfaces": interfaces,
        "functions": functions
    }

def main():
    root_dir = "."
    map_content = "# ThePress Bot Ultimate - Ամբողջական Ճարտարապետական Քարտեզ\n\n"
    map_content += "Այս փաստաթուղթը պարունակում է բոտի ամբողջական կոդի քարտեզը, ներառյալ ֆայլերը, ֆունկցիաները, կառուցվածքները (structs) և տվյալների բազայի սխեման:\n\n"

    map_content += "## Տվյալների Բազա (Database)\n"
    map_content += "Բազան SQLite է (`bot_ultimate.db`): Հիմնական աղյուսակներն են.\n"
    map_content += "- `app_settings` - բոտի գլոբալ կարգավորումներ (WP, AI keys, prompts)\n"
    map_content += "- `rss_topics` - RSS հոսքերի ցանկ\n"
    map_content += "- `articles` - հավաքագրված և մշակված հոդվածների պահոց (կարգավիճակներով՝ pending, published, failed)\n\n"

    map_content += "## Ֆայլերի և Ֆունկցիաների Ցանկ\n\n"

    for subdir, _, files in os.walk(root_dir):
        if '.git' in subdir: continue
        for file in sorted(files):
            if file.endswith('.go'):
                filepath = os.path.join(subdir, file)
                info = parse_go_file(filepath)
                rel_path = os.path.relpath(filepath, root_dir)
                map_content += f"### `{rel_path}` (Package: `{info['package']}`)\n"
                if info['structs']:
                    map_content += "**Կառուցվածքներ (Structs):** " + ", ".join(info['structs']) + "\n"
                if info['interfaces']:
                    map_content += "**Ինտերֆեյսներ (Interfaces):** " + ", ".join(info['interfaces']) + "\n"
                if info['functions']:
                    map_content += "**Ֆունկցիաներ:**\n"
                    for f in info['functions']:
                        map_content += f"- `{f}`\n"
                map_content += "\n"

    with open(os.path.join(root_dir, "CODEBASE_MAP_AR.md"), "w", encoding="utf-8") as f:
        f.write(map_content)

if __name__ == "__main__":
    main()

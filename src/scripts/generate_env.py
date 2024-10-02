import configparser

def parse_config_to_env(config_file, env_file):
    config = configparser.ConfigParser()
    config.read(config_file)

    def evaluate_value(value):
        try:
            # Try to evaluate expressions like "7*24*60*60"
            evaluated_value = str(eval(value))
            return evaluated_value
        except Exception:
            # Return the original value if it's not a valid expression
            return value

    with open(env_file, 'w') as env:
        for section in config.sections():
            # Replace section headers to environment-friendly format
            section_name = section.replace('.', '_').upper()

            # Write section header comment
            env.write(f"# === {section_name} ===\n")

            for key, value in config.items(section):
                # Convert key to uppercase and join it with section name
                env_var_name = f"{section_name}_{key}".upper()

                # Evaluate the value if possible
                evaluated_value = evaluate_value(value)

                # Write variable in the form VAR=VALUE
                env.write(f"{env_var_name}={evaluated_value}\n")
            
            env.write("\n")

if __name__ == "__main__":
    parse_config_to_env("config.toml", ".env")

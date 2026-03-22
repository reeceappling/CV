# Execute commands in a subshell so that the final directory is the initial
(
  current_dir_name=$(basename "$PWD")
  case "$current_dir_name" in
    "CV")
      echo "working from CV directory"
      ;;
    "scripts")
      echo "working from scripts directory"
      cd ../ || exit 1
      ;;
    *)
      echo "unhandled working directory: $current_dir_name"
      exit 1
      ;;
  esac
  chmod -R 777 ./quartz
  go run ./generator
)

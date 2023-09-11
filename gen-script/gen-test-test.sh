#!/bin/bash

indiog_dir="/Volumes/MicronSSD/Users/henoya/src/Bluesky_repo/indigo"

# indigo の struct 定義を変換する
read -d '' scriptVariable << 'EOF'
BEGIN{
  import_block=0;
  struct_block=0;
}
{
  if (import_block == 0) {
    match($0, /import \\(/, r);
    if (length(r) >= 1) {
      print $0;
      import_block=1;
    }
  } else {
    if (import_block == 1) {
      match($0, /\\)/, r);
      if (length(r) >= 1) {
        import_block=2;
        print $0;
      } else {
        print $0;
      }
    }
  }
  if (import_block > 1) {
    if (struct_block == 0) {
      match($0, /type +(.+) +struct +{/, r);
      if (length(r) >= 1) {
        struct_block=1;
        print $0;
      }
    } else {
      match($0, /}/, r);
      if (length(r) >= 1) {
        struct_block=0;
        print $0;
      } else {
        match($0, /^(.+)(.json:.+)$/, r);
        # print length(r);
        # for (i=1; i<=length(r); i++) {
        #   print i,":",r[i];
        # }
        if (length(r) >= 1) {
          print r[1];
        } else {
          print $0;
        }
      }
    }
  }
}
EOF

if [[ -f "${indiog_dir}/$1" ]]; then
  file_base_name=$(basename $1)
  target_file="./${file_base_name}"
  echo "// Code generated fron $1 DO NOT EDIT."
  echo "package ${GOPACKAGE}"
  gawk "${scriptVariable}" "${indiog_dir}/$1"
fi

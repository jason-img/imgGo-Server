import datetime
import os
import re
import subprocess

import py7zr
import yaml


def back_upload_files(work_path: str, back_path: str, back_keep_days: int) -> None:
    # 执行备份
    back_file = f'{os.path.join(back_path, datetime.datetime.now().strftime("%Y-%m-%d"))}.7z'
    print('back_file ->', back_file)
    with py7zr.SevenZipFile(back_file, 'w') as archive:
        archive.writeall(work_path)

    # 删除修改日期N天前的备份
    command = r'sudo find %s -maxdepth 1 -type f -mtime +%d -exec rm -f {} \;' % (back_path, back_keep_days)
    print('command ->', command)
    status, output = subprocess.getstatusoutput(command)
    print('status ->', status)
    print('output ->', output)

    print('Done.')


# 获取环境变量
def get_env() -> str:
    try:
        file = open('./.env', 'r', encoding='utf-8')
        data = file.readline()
        file.close()
        data = re.match(r'\S*', data, flags=0)
        return data.group()
    except Exception as e:
        print(e)
        print('获取环境变量失败')
        exit(-1)


def get_path(yaml_file: str) -> (str, str, int):
    try:
        file = open(yaml_file, 'r', encoding="utf-8")
        file_data = file.read()
        file.close()

        cfg = yaml.load(file_data, Loader=yaml.Loader)

        env = get_env()

        cfg = cfg['PrdEnv'] if env.lower() == 'prd' else cfg['DevEnv']

        return cfg['Path'], cfg['BackupPath'], cfg['BackupKeepDays']
    except Exception as e:
        print(e)
        print('获取配置失败')
        exit(-1)


if __name__ == '__main__':
    back_upload_files(*get_path('./config.yml'))

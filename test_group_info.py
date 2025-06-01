#!/usr/bin/env python3
"""
测试群信息接口的脚本
用于验证 get_group_info 接口是否正常工作
"""

import requests
import json

def test_group_info(group_id="978880814"):
    """测试获取群信息接口"""
    url = "http://127.0.0.1:3000/get_group_info"
    
    payload = json.dumps({
        "group_id": group_id,
        "no_cache": True
    })
    
    headers = {
        'Content-Type': 'application/json'
    }
    
    try:
        response = requests.request("POST", url, headers=headers, data=payload)
        print(f"状态码: {response.status_code}")
        print(f"响应内容: {response.text}")
        
        if response.status_code == 200:
            data = response.json()
            if data.get("status") == "ok" and data.get("retcode") == 0:
                group_data = data.get("data", {})
                member_count = group_data.get("member_count", 0)
                max_member_count = group_data.get("max_member_count", 0)
                
                print(f"\n群信息解析:")
                print(f"群号: {group_data.get('group_id')}")
                print(f"群名: {group_data.get('group_name')}")
                print(f"当前人数: {member_count}")
                print(f"最大人数: {max_member_count}")
                
                if member_count >= max_member_count:
                    print("状态: 群聊已满 ❌")
                    print("新的加群申请将被拒绝")
                else:
                    remaining = max_member_count - member_count
                    print(f"状态: 群聊未满 ✅ (还可加入{remaining}人)")
                
                return True
            else:
                print(f"接口返回错误: {data}")
                return False
        else:
            print(f"HTTP请求失败: {response.status_code}")
            return False
            
    except Exception as e:
        print(f"请求异常: {e}")
        return False

if __name__ == "__main__":
    print("测试群信息接口...")
    print("=" * 50)
    
    # 测试指定群号
    test_group_info("978880814")
    
    print("\n" + "=" * 50)
    print("测试完成")
    print("\n使用说明:")
    print("1. 确保 NapCat 正在运行并监听 3000 端口")
    print("2. 确保机器人已加入要测试的群聊")
    print("3. 在配置文件中启用 enable-reject-full-group-join-request: true")
    print("4. 当群聊满员时，新的加群申请将被自动拒绝")

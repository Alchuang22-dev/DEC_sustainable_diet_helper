import torch
from torchvision import datasets, transforms
import torch.nn as nn
from model import SimpleNet  # 假设你的模型在 model.py 中

def evaluate_model():
    # 设置设备（GPU/CPU）
    device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
    
    # 数据预处理，与训练时保持一致
    transform = transforms.Compose([
        transforms.Resize((256, 256)),  # 与训练时相同的 Resize 操作
        transforms.ToTensor(),
        transforms.Normalize((0.5, 0.5, 0.5), (0.5, 0.5, 0.5))
    ])
    
    # 加载测试集
    testset = datasets.Food101(root='./data', split='test', download=True, transform=transform)
    testloader = torch.utils.data.DataLoader(testset, batch_size=32, shuffle=False)
    
    # 加载训练好的模型，并确保类别数一致（这里假设模型最后输出层已经修改为101个类别）
    model = SimpleNet()  
    # 如有保存的模型权重，可以加载：
    # model.load_state_dict(torch.load('model.pth'))
    model.to(device)
    model.eval()  # 设定为评估模式

    criterion = nn.CrossEntropyLoss()
    
    total = 0
    correct = 0
    test_loss = 0.0
    
    # 关闭梯度计算，提高评估速度
    with torch.no_grad():
        for inputs, labels in testloader:
            inputs, labels = inputs.to(device), labels.to(device)
            
            outputs = model(inputs)
            loss = criterion(outputs, labels)
            test_loss += loss.item() * inputs.size(0)
            
            # 获取预测结果
            _, predicted = torch.max(outputs, 1)
            total += labels.size(0)
            correct += (predicted == labels).sum().item()
    
    avg_loss = test_loss / total
    accuracy = correct / total
    
    print(f"Test Loss: {avg_loss:.4f}")
    print(f"Test Accuracy: {accuracy:.4f}")

if __name__ == '__main__':
    evaluate_model()

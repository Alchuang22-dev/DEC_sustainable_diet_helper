import torch
import torch.nn as nn
from torchvision import models, transforms
from PIL import Image
import matplotlib.pyplot as plt
import numpy as np
import time

def create_model(model_name='resnet50', num_classes=101):
    """创建模型架构"""
    if model_name == 'resnet18':
        model = models.resnet18(weights=None)
    elif model_name == 'resnet50':
        model = models.resnet50(weights=None)
    elif model_name == 'resnet101':
        model = models.resnet101(weights=None)
    else:
        raise ValueError(f"不支持的模型: {model_name}")
    
    # 修改最后的全连接层以适应Food-101的101个类别
    num_features = model.fc.in_features
    model.fc = nn.Sequential(
        nn.Dropout(0.5),
        nn.Linear(num_features, num_classes)
    )
    
    return model

def load_model(model_path, model_name='resnet50', num_classes=101):
    """加载保存的模型"""
    # 创建模型架构
    model = create_model(model_name, num_classes)
    
    # 加载保存的模型权重和类别名称
    checkpoint = torch.load(model_path, map_location=torch.device('cpu'))
    
    model.load_state_dict(checkpoint['model_state_dict'])
    class_names = checkpoint['class_names']
    
    # 设置为评估模式
    model.eval()
    
    return model, class_names

def preprocess_image(image_path):
    """预处理图像以符合模型输入要求"""
    # 使用与训练时相同的预处理步骤
    transform = transforms.Compose([
        transforms.Resize(256),
        transforms.CenterCrop(224),
        transforms.ToTensor(),
        transforms.Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])
    ])
    
    # 打开图像并应用转换
    image = Image.open(image_path).convert('RGB')
    image_tensor = transform(image)
    
    # 添加批次维度
    image_tensor = image_tensor.unsqueeze(0)
    
    return image_tensor, image

def predict_image(model, image_tensor, class_names, device='cpu', top_k=5):
    """使用模型预测图像类别"""
    # 移动数据到指定设备
    image_tensor = image_tensor.to(device)
    model = model.to(device)
    
    # 关闭梯度计算
    with torch.no_grad():
        outputs = model(image_tensor)
        probabilities = torch.nn.functional.softmax(outputs, dim=1)
        
        # 获取top-k预测结果
        top_probs, top_indices = torch.topk(probabilities, top_k)
        
    # 转换为numpy数组
    top_probs = top_probs.cpu().numpy().flatten()
    top_indices = top_indices.cpu().numpy().flatten()
    
    # 获取对应的类别名称
    top_classes = [class_names[idx] for idx in top_indices]
    
    return top_probs, top_classes

def display_results(image, top_probs, top_classes, 
                   show_plot=False, save_plot=True, 
                   save_path='prediction_result.png',
                   print_results=True):
    """灵活处理结果显示"""
    
    if print_results:
        print("预测结果:")
        for i, (cls, prob) in enumerate(zip(top_classes, top_probs)):
            print(f"{i+1}. {cls}: {prob*100:.2f}%")
    
    # 创建可视化
    plt.figure(figsize=(12, 6))
    
    # 显示图像
    plt.subplot(1, 2, 1)
    plt.imshow(image)
    plt.axis('off')
    plt.title('输入图像')
    
    # 显示预测结果
    plt.subplot(1, 2, 2)
    y_pos = np.arange(len(top_classes))
    plt.barh(y_pos, top_probs * 100)
    plt.yticks(y_pos, top_classes)
    plt.xlabel('概率 (%)')
    plt.title('预测结果')
    # plt.tight_layout()
    
    # 保存图像
    if save_plot:
        plt.savefig(save_path)
        print(f"预测结果已保存到 {save_path}")
    
    # 尝试显示图像
    if show_plot:
        try:
            plt.show()
        except Exception as e:
            print(f"警告: 无法显示图像 ({e}). 结果已保存到文件.")
    
    plt.close()

# 主函数
def classify_custom_image(model_path, image_path, model_name='resnet50', device='cpu', top_k=5):
    """对自定义图像进行分类"""
    # 加载模型和类别名称
    model, class_names = load_model(model_path, model_name)
    
    # 预处理图像
    image_tensor, original_image = preprocess_image(image_path)
    
    # 预测
    top_probs, top_classes = predict_image(model, image_tensor, class_names, device, top_k)
    
    # 显示结果
    display_results(original_image, top_probs, top_classes)
    
    return top_probs, top_classes

# 使用示例
if __name__ == "__main__":
    # 计时器
    start = time.perf_counter()
    
    # 模型路径
    model_path = 'food101_model_complete.pth'
    
    # 图像路径
    image_path = 'food_image.jpg'
    
    # 设置设备
    device = 'cuda' if torch.cuda.is_available() else 'cpu'
    
    # 分类图像
    top_probs, top_classes = classify_custom_image(
        model_path=model_path,
        image_path=image_path,
        model_name='resnet50',  # 使用的模型架构
        device=device,
        top_k=5  # 显示前5个预测结果
    )
    
    # 计时器
    end = time.perf_counter()
    runTime = end - start
    print("运行时间：", runTime)
    
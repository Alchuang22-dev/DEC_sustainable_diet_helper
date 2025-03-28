import torch
import torch.nn as nn
from torchvision import models, transforms
from PIL import Image
import matplotlib.pyplot as plt
import numpy as np
import time

class EfficientBinaryClassifier(nn.Module):
    def __init__(self, base_model_name='resnet50'):
        """
        使用单一骨干网络的二分类器
        
        参数:
            base_model_name: 基础模型架构
        """
        super(EfficientBinaryClassifier, self).__init__()
        
        # 创建基础模型（使用ImageNet预训练权重）
        if base_model_name == 'resnet18':
            base_model = models.resnet18(weights='IMAGENET1K_V1')
            num_features = base_model.fc.in_features
        elif base_model_name == 'resnet50':
            base_model = models.resnet50(weights='IMAGENET1K_V1')
            num_features = base_model.fc.in_features
        else:
            raise ValueError(f"不支持的模型: {base_model_name}")
        
        # 去掉分类层
        self.features = nn.Sequential(*list(base_model.children())[:-1])
        
        # 二分类器
        self.classifier = nn.Sequential(
            nn.Flatten(),
            nn.Linear(num_features, 512),
            nn.ReLU(),
            nn.Dropout(0.5),
            nn.Linear(512, 2)
        )
    
    def forward(self, x):
        """前向传播"""
        # 特征提取
        features = self.features(x)
        
        # 分类
        output = self.classifier(features)
        
        return output

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
    
    # 修改最后的全连接层
    num_features = model.fc.in_features
    model.fc = nn.Sequential(
        nn.Dropout(0.5),
        nn.Linear(num_features, num_classes)
    )
    
    return model

def load_model(model_path, model_name='resnet50', num_classes=101):
    """加载模型"""
    model = create_model(model_name, num_classes)
    checkpoint = torch.load(model_path, map_location='cpu')
    
    if 'model_state_dict' in checkpoint:
        model.load_state_dict(checkpoint['model_state_dict'])
        class_names = checkpoint.get('class_names', [f"Class_{i}" for i in range(num_classes)])
    else:
        model.load_state_dict(checkpoint)
        class_names = [f"Class_{i}" for i in range(num_classes)]
    
    model.eval()
    return model, class_names

def load_binary_classifier(model_path, device='cpu'):
    """加载二分类模型"""
    model = EfficientBinaryClassifier(base_model_name='resnet50')
    checkpoint = torch.load(model_path, map_location=device)
    model.load_state_dict(checkpoint['model_state_dict'])
    model = model.to(device)
    model.eval()
    return model

def predict_food_category(image_path, binary_model_path, food101_model_path, custom_model_path, device='cpu'):
    """
    使用二分类器和两个细分类模型进行预测
    
    参数:
        image_path: 图像路径
        binary_model_path: 二分类模型路径
        food101_model_path: Food-101模型路径
        custom_model_path: 本地数据库模型路径
        device: 计算设备
    
    返回:
        预测结果
    """
    # 图像预处理
    transform = transforms.Compose([
        transforms.Resize(256),
        transforms.CenterCrop(224),
        transforms.ToTensor(),
        transforms.Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])
    ])
    
    # 加载图像
    image = Image.open(image_path).convert('RGB')
    image_tensor = transform(image).unsqueeze(0).to(device)
    
    # 加载二分类模型
    binary_classifier = load_binary_classifier(binary_model_path, device)
    
    # 预测数据集类别
    with torch.no_grad():
        binary_outputs = binary_classifier(image_tensor)
        binary_probs = torch.nn.functional.softmax(binary_outputs, dim=1)
        dataset_idx = torch.argmax(binary_probs, dim=1).item()
        dataset_prob = binary_probs[0, dataset_idx].item()
    
    # 根据数据集类别选择相应的细分类模型
    if dataset_idx == 0:  # Food-101
        model, class_names = load_model(food101_model_path, 'resnet50', 101)
        dataset_name = "Food-101"
    else:  # 本地数据库
        model, class_names = load_model(custom_model_path, 'resnet50', 208)
        dataset_name = "本地数据库"
    
    model = model.to(device)
    
    # 预测细分类别
    with torch.no_grad():
        outputs = model(image_tensor)
        probabilities = torch.nn.functional.softmax(outputs, dim=1)
        
        # 获取前5个预测结果
        top5_prob, top5_idx = torch.topk(probabilities, min(5, len(class_names)))
    
    # 输出预测结果
    print(f"预测结果 - 图像: {image_path}")
    print(f"数据集: {dataset_name} (置信度: {dataset_prob*100:.2f}%)")
    print("类别预测:")
    
    results = []
    for i in range(min(5, len(top5_idx[0]))):
        idx = top5_idx[0, i].item()
        prob = top5_prob[0, i].item()
        class_name = class_names[idx]
        print(f"{i+1}. {class_name}: {prob*100:.2f}%")
        results.append((class_name, prob))
    
    return dataset_name, dataset_prob, results

# 使用示例
if __name__ == "__main__":
     # 计时器
    start = time.perf_counter()
    
    # 设置设备
    device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
    
    # 模型路径（模型路径在model-v2/-v3/classified中已经设计好了）
    binary_model_path = 'efficient_binary_classifier.pth'
    food101_model_path = 'food101_model_complete.pth'
    custom_model_path = 'custom_food_model_complete.pth'
    
    # 测试图像
    image_path = 'food_image_2.jpg'
    
    # 预测
    predict_food_category(
        image_path=image_path,
        binary_model_path=binary_model_path,
        food101_model_path=food101_model_path,
        custom_model_path=custom_model_path,
        device=device
    )
    
    # 计时器
    end = time.perf_counter()
    runTime = end - start
    print("运行时间：", runTime)
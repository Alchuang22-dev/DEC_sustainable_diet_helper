import os
import torch
import torch.nn as nn
from torch.utils.data import Dataset, DataLoader
from torchvision import models, transforms
from PIL import Image
import matplotlib.pyplot as plt
import numpy as np
import time
import copy

# 自定义数据集类
class CustomFoodDataset(Dataset):
    def __init__(self, root_dir, list_file, name_file, transform=None):
        """
        自定义食物数据集
        
        参数:
            root_dir (string): 数据集根目录
            list_file (string): 包含图片路径和标签的文件
            name_file (string): 包含类别名称的文件
            transform (callable, optional): 图像转换操作
        """
        self.root_dir = root_dir
        self.transform = transform
        
        # 加载类别名称
        self.class_names = []
        self.class_dict = {}  # 用于存储类别ID到索引的映射
        
        with open(os.path.join(root_dir, name_file), 'r', encoding='utf-8') as f:
            for line in f:
                parts = line.strip().split('\t')
                if len(parts) >= 3:
                    class_id = int(parts[0])
                    chinese_name = parts[1]
                    english_name = parts[2]
                    self.class_names.append(f"{chinese_name} ({english_name})")
                    self.class_dict[class_id] = len(self.class_names) - 1
        
        # 加载图片路径和标签
        self.image_paths = []
        self.labels = []
        
        split_type = list_file.split('_')[0]  # 获取train, val或test
        
        with open(os.path.join(root_dir, list_file), 'r') as f:
            for line in f:
                parts = line.strip().split()
                if len(parts) >= 2:
                    img_path = parts[0]
                    label = int(parts[1])
                    
                    # 构建完整的图像路径
                    full_img_path = os.path.join(root_dir, split_type, img_path)
                    
                    self.image_paths.append(full_img_path)
                    self.labels.append(self.class_dict[label])  # 使用映射后的索引
    
    def __len__(self):
        return len(self.image_paths)
    
    def __getitem__(self, idx):
        img_path = self.image_paths[idx]
        label = self.labels[idx]
        
        # 读取图像
        try:
            image = Image.open(img_path).convert('RGB')
        except Exception as e:
            print(f"无法加载图像 {img_path}: {e}")
            # 如果图像加载失败，返回黑色图像
            image = Image.new('RGB', (224, 224))
        
        if self.transform:
            image = self.transform(image)
            
        return image, label

# 创建和加载模型函数
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

def prepare_datasets(root_dir='data/dataset_release/release_data'):
    """
    准备训练、验证和测试数据集
    
    参数:
        root_dir: 数据集根目录
    返回:
        dataloaders: 包含train, val和test的DataLoader字典
        dataset_sizes: 各数据集的大小
        class_names: 类别名称
    """
    # 定义数据转换
    data_transforms = {
        'train': transforms.Compose([
            transforms.RandomResizedCrop(224),
            transforms.RandomHorizontalFlip(),
            transforms.ColorJitter(brightness=0.2, contrast=0.2, saturation=0.2),
            transforms.ToTensor(),
            transforms.Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])
        ]),
        'val': transforms.Compose([
            transforms.Resize(256),
            transforms.CenterCrop(224),
            transforms.ToTensor(),
            transforms.Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])
        ]),
        'test': transforms.Compose([
            transforms.Resize(256),
            transforms.CenterCrop(224),
            transforms.ToTensor(),
            transforms.Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])
        ]),
    }
    
    # 创建数据集
    train_dataset = CustomFoodDataset(
        root_dir=root_dir,
        list_file='train_list.txt',
        name_file='name.txt',
        transform=data_transforms['train']
    )
    
    val_dataset = CustomFoodDataset(
        root_dir=root_dir,
        list_file='val_list.txt',
        name_file='name.txt',
        transform=data_transforms['val']
    )
    
    test_dataset = CustomFoodDataset(
        root_dir=root_dir,
        list_file='test_list.txt',
        name_file='name.txt',
        transform=data_transforms['test']
    )
    
    # 创建数据加载器
    batch_size = 32
    dataloaders = {
        'train': DataLoader(train_dataset, batch_size=batch_size, shuffle=True, num_workers=4),
        'val': DataLoader(val_dataset, batch_size=batch_size, shuffle=False, num_workers=4),
        'test': DataLoader(test_dataset, batch_size=batch_size, shuffle=False, num_workers=4)
    }
    
    dataset_sizes = {
        'train': len(train_dataset),
        'val': len(val_dataset),
        'test': len(test_dataset)
    }
    
    return dataloaders, dataset_sizes, train_dataset.class_names

def adapt_model_for_finetuning(model_path, model_name='resnet50', source_classes=101, target_classes=208):
    """
    调整预训练模型以适应新的分类任务
    
    参数:
        model_path: 预训练模型的路径
        model_name: 模型架构名称
        source_classes: 原始模型的类别数
        target_classes: 目标数据集的类别数
    """
    # 加载预训练模型
    model = create_model(model_name, source_classes)
    checkpoint = torch.load(model_path, map_location=torch.device('cpu'))
    model.load_state_dict(checkpoint['model_state_dict'])
    
    # 获取特征提取部分的输入特征数
    num_features = model.fc[1].in_features
    
    # 替换分类器层
    model.fc = nn.Sequential(
        nn.Dropout(0.5),
        nn.Linear(num_features, target_classes)
    )
    
    return model

def finetune_model(model, dataloaders, dataset_sizes, criterion, optimizer, scheduler, num_epochs=25, device='cpu'):
    """
    微调模型
    
    参数:
        model: 需要微调的模型
        dataloaders: 包含train和val的DataLoader
        dataset_sizes: 数据集大小
        criterion: 损失函数
        optimizer: 优化器
        scheduler: 学习率调度器
        num_epochs: 训练轮数
        device: 计算设备
    """
    since = time.time()
    
    best_model_wts = model.state_dict()
    best_acc = 0.0
    
    for epoch in range(num_epochs):
        print(f'Epoch {epoch+1}/{num_epochs}')
        print('-' * 10)
        
        # 每个epoch有训练和验证阶段
        for phase in ['train', 'val']:
            if phase == 'train':
                model.train()  # 设置模型为训练模式
            else:
                model.eval()   # 设置模型为评估模式
            
            running_loss = 0.0
            running_corrects = 0
            
            # 迭代数据
            for inputs, labels in dataloaders[phase]:
                inputs = inputs.to(device)
                labels = labels.to(device)
                
                # 梯度清零
                optimizer.zero_grad()
                
                # 前向传播
                with torch.set_grad_enabled(phase == 'train'):
                    outputs = model(inputs)
                    _, preds = torch.max(outputs, 1)
                    loss = criterion(outputs, labels)
                    
                    # 在训练阶段，反向传播和优化
                    if phase == 'train':
                        loss.backward()
                        optimizer.step()
                
                # 统计
                running_loss += loss.item() * inputs.size(0)
                running_corrects += torch.sum(preds == labels.data)
            
            if phase == 'train' and scheduler is not None and not isinstance(scheduler, torch.optim.lr_scheduler.ReduceLROnPlateau):
                scheduler.step()
            
            epoch_loss = running_loss / dataset_sizes[phase]
            epoch_acc = running_corrects.double() / dataset_sizes[phase]
            
            print(f'{phase} Loss: {epoch_loss:.4f} Acc: {epoch_acc:.4f}')
            
            # 如果是验证阶段并且有更好的准确率，保存模型
            if phase == 'val' and epoch_acc > best_acc:
                best_acc = epoch_acc
                best_model_wts = model.state_dict().copy()
                # 保存最佳模型
                torch.save({
                    'model_state_dict': best_model_wts,
                    'optimizer_state_dict': optimizer.state_dict(),
                    'epoch': epoch,
                    'accuracy': best_acc
                }, f'best_model_{num_epochs}epochs.pth')
                print(f'保存新的最佳模型，精度: {best_acc:.4f}')
            
            # 如果使用ReduceLROnPlateau，在验证阶段后调整学习率
            if phase == 'val' and scheduler is not None and isinstance(scheduler, torch.optim.lr_scheduler.ReduceLROnPlateau):
                scheduler.step(epoch_loss)
        
        print()
    
    time_elapsed = time.time() - since
    print(f'训练完成，用时 {time_elapsed // 60:.0f}m {time_elapsed % 60:.0f}s')
    print(f'最佳验证精度: {best_acc:.4f}')
    
    # 加载最佳模型权重
    model.load_state_dict(best_model_wts)
    return model

def main():
    # 设置随机种子以确保可重复性
    torch.manual_seed(42)
    np.random.seed(42)
    
    # 设置设备
    device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
    print(f"使用设备: {device}")
    
    # 准备数据集
    root_dir = 'data/dataset_release/release_data'
    dataloaders, dataset_sizes, class_names = prepare_datasets(root_dir)
    print(f"数据集大小: 训练={dataset_sizes['train']}, 验证={dataset_sizes['val']}, 测试={dataset_sizes['test']}")
    print(f"类别数量: {len(class_names)}")
    
    # 加载预训练模型并调整分类器
    model_path = 'food101_model_complete.pth'
    model = adapt_model_for_finetuning(
        model_path=model_path,
        model_name='resnet50',
        source_classes=101,
        target_classes=208  # 根据您的数据集类别数量调整
    )
    
    # 将模型移动到设备
    model = model.to(device)
    
    # 设置损失函数和优化器
    criterion = nn.CrossEntropyLoss()
    
    # 微调策略：冻结大部分层，只训练新的分类器和最后几层
    # 1. 首先冻结所有参数
    for param in model.parameters():
        param.requires_grad = False
    
    # 2. 解冻最后几层和分类器
    for param in model.layer4.parameters():  # 解冻最后一个残差块
        param.requires_grad = True
    for param in model.fc.parameters():  # 解冻分类器
        param.requires_grad = True
    
    # 3. 设置不同的学习率
    params_to_update = [
        {'params': model.layer4.parameters(), 'lr': 0.0001},
        {'params': model.fc.parameters(), 'lr': 0.001}
    ]
    
    optimizer = torch.optim.Adam(params_to_update)
    
    # 学习率调度器
    scheduler = torch.optim.lr_scheduler.ReduceLROnPlateau(
        optimizer, mode='min', factor=0.1, patience=3, verbose=True
    )
    
    print("step 1")
    # 进行第一阶段微调（只训练分类器和最后几层）
    model = finetune_model(
        model=model,
        dataloaders=dataloaders,
        dataset_sizes=dataset_sizes,
        criterion=criterion,
        optimizer=optimizer,
        scheduler=scheduler,
        num_epochs=10,
        device=device
    )
    
    print("step 2")
    # 第二阶段：解冻更多层，使用较小的学习率进行全面微调
    for param in model.parameters():
        param.requires_grad = True
    
    # 设置不同的学习率
    params_to_update = [
        {'params': model.conv1.parameters(), 'lr': 1e-5},
        {'params': model.bn1.parameters(), 'lr': 1e-5},
        {'params': model.layer1.parameters(), 'lr': 1e-5},
        {'params': model.layer2.parameters(), 'lr': 1e-5},
        {'params': model.layer3.parameters(), 'lr': 1e-4},
        {'params': model.layer4.parameters(), 'lr': 1e-4},
        {'params': model.fc.parameters(), 'lr': 1e-3}
    ]
    
    optimizer = torch.optim.Adam(params_to_update)
    
    # 学习率调度器
    scheduler = torch.optim.lr_scheduler.ReduceLROnPlateau(
        optimizer, mode='min', factor=0.1, patience=3, verbose=True
    )
    
    # 进行第二阶段微调
    model = finetune_model(
        model=model,
        dataloaders=dataloaders,
        dataset_sizes=dataset_sizes,
        criterion=criterion,
        optimizer=optimizer,
        scheduler=scheduler,
        num_epochs=15,
        device=device
    )
    
    # 保存完整模型
    torch.save({
        'model_state_dict': model.state_dict(),
        'optimizer_state_dict': optimizer.state_dict(),
        'class_names': class_names
    }, 'custom_food_model_complete.pth')
    print("模型已保存为 'custom_food_model_complete.pth'")
    
    # 评估最终模型
    print("在测试集上评估最终模型...")
    model.eval()
    running_corrects = 0
    total = 0
    
    with torch.no_grad():
        for inputs, labels in dataloaders['test']:
            inputs = inputs.to(device)
            labels = labels.to(device)
            
            outputs = model(inputs)
            _, preds = torch.max(outputs, 1)
            
            running_corrects += torch.sum(preds == labels.data)
            total += labels.size(0)
    
    test_acc = running_corrects.double() / total
    print(f'最终测试集准确率: {test_acc:.4f}')

if __name__ == "__main__":
    main()
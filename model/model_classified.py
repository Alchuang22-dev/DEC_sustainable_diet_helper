import os
import torch
import torch.nn as nn
from torch.utils.data import Dataset, DataLoader, ConcatDataset
from torchvision import models, transforms
from PIL import Image
import time
import copy
import numpy as np

# 为Food-101数据集创建自定义数据集类
class Food101Dataset(Dataset):
    def __init__(self, data_dir, transform=None, split='train'):
        """
        Food-101数据集加载器
        
        参数:
            data_dir: Food-101数据集根目录
            transform: 图像转换
            split: 'train'或'test'
        """
        self.data_dir = data_dir
        self.transform = transform
        self.split = split
        
        # 加载类别名称
        self.class_names = []
        with open(os.path.join(data_dir, 'meta', 'classes.txt'), 'r') as f:
            for line in f:
                self.class_names.append(line.strip())
        
        # 创建类别到索引的映射
        self.class_to_idx = {cls: i for i, cls in enumerate(self.class_names)}
        
        # 加载图像路径和标签
        self.image_paths = []
        self.labels = []
        
        with open(os.path.join(data_dir, 'meta', f'{split}.txt'), 'r') as f:
            for line in f:
                path = line.strip()
                cls = path.split('/')[0]
                img_path = os.path.join(data_dir, 'images', path + '.jpg')
                
                self.image_paths.append(img_path)
                self.labels.append(self.class_to_idx[cls])
        
        print(f"Food-101 {split}: 加载了 {len(self.image_paths)} 张图像，共 {len(self.class_names)} 个类别")
    
    def __len__(self):
        return len(self.image_paths)
    
    def __getitem__(self, idx):
        img_path = self.image_paths[idx]
        label = self.labels[idx]
        
        try:
            image = Image.open(img_path).convert('RGB')
        except Exception as e:
            print(f"无法加载图像 {img_path}: {e}")
            image = Image.new('RGB', (224, 224))
        
        if self.transform:
            image = self.transform(image)
            
        return image, 0  # 0表示来自Food-101数据集

# 为本地数据库创建自定义数据集类
class CustomFoodDataset(Dataset):
    def __init__(self, root_dir, list_file, name_file, transform=None):
        """
        自定义食物数据集
        
        参数:
            root_dir: 数据集根目录
            list_file: 包含图片路径和标签的文件
            name_file: 包含类别名称的文件
            transform: 图像转换操作
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
                    
                    if label in self.class_dict:
                        self.image_paths.append(full_img_path)
                        self.labels.append(self.class_dict[label])
                    else:
                        print(f"警告: 忽略标签 {label}，因为它不在类别映射中")
        
        print(f"本地数据库 {split_type}: 加载了 {len(self.image_paths)} 张图像，共 {len(self.class_names)} 个类别")
    
    def __len__(self):
        return len(self.image_paths)
    
    def __getitem__(self, idx):
        img_path = self.image_paths[idx]
        label = self.labels[idx]
        
        try:
            image = Image.open(img_path).convert('RGB')
        except Exception as e:
            print(f"无法加载图像 {img_path}: {e}")
            image = Image.new('RGB', (224, 224))
        
        if self.transform:
            image = self.transform(image)
            
        return image, 1  # 1表示来自本地数据集

# 特征提取器类
class FeatureExtractor(nn.Module):
    def __init__(self, model_path, model_name='resnet50', num_classes=101):
        """
        从预训练模型创建特征提取器
        
        参数:
            model_path: 预训练模型的路径
            model_name: 模型架构
            num_classes: 原模型的类别数
        """
        super(FeatureExtractor, self).__init__()
        
        # 创建基础模型
        if model_name == 'resnet18':
            self.model = models.resnet18(weights=None)
            num_features = self.model.fc.in_features
            self.model.fc = nn.Linear(num_features, num_classes)
        elif model_name == 'resnet50':
            self.model = models.resnet50(weights=None)
            num_features = self.model.fc.in_features
            self.model.fc = nn.Linear(num_features, num_classes)
        elif model_name == 'resnet101':
            self.model = models.resnet101(weights=None)
            num_features = self.model.fc.in_features
            self.model.fc = nn.Linear(num_features, num_classes)
        else:
            raise ValueError(f"不支持的模型: {model_name}")
        
        # 加载预训练权重
        checkpoint = torch.load(model_path, map_location='cpu')
        if 'model_state_dict' in checkpoint:
            self.model.load_state_dict(checkpoint['model_state_dict'])
        else:
            self.model.load_state_dict(checkpoint)
        
        # 提取特征的模型（去掉分类层）
        self.features = nn.Sequential(*list(self.model.children())[:-1])
        
        # 冻结所有参数
        for param in self.features.parameters():
            param.requires_grad = False
    
    def forward(self, x):
        """提取特征"""
        x = self.features(x)
        x = torch.flatten(x, 1)
        return x

# 二分类器模型
class BinaryClassifier(nn.Module):
    def __init__(self, food101_model_path, custom_model_path, feature_dim=2048):
        """
        使用两个预训练模型作为特征提取器的二分类模型
        
        参数:
            food101_model_path: Food-101模型路径
            custom_model_path: 本地数据库模型路径
            feature_dim: 特征维度
        """
        super(BinaryClassifier, self).__init__()
        
        # 创建特征提取器
        self.food101_extractor = FeatureExtractor(
            model_path=food101_model_path,
            model_name='resnet50',
            num_classes=101
        )
        
        self.custom_extractor = FeatureExtractor(
            model_path=custom_model_path,
            model_name='resnet50',
            num_classes=208  # 假设本地数据库有208个类别
        )
        
        # 二分类器
        self.classifier = nn.Sequential(
            nn.Linear(feature_dim * 2, 512),  # 连接两个特征提取器的输出
            nn.ReLU(),
            nn.Dropout(0.5),
            nn.Linear(512, 2)
        )
    
    def forward(self, x):
        """前向传播"""
        # 从两个特征提取器获取特征
        food101_features = self.food101_extractor(x)
        custom_features = self.custom_extractor(x)
        
        # 连接特征
        combined_features = torch.cat((food101_features, custom_features), dim=1)
        
        # 分类
        output = self.classifier(combined_features)
        
        return output

# 准备数据集
def prepare_binary_datasets(food101_dir, local_data_dir):
    """准备二分类数据集"""
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
    
    # 创建Food-101数据集
    food101_train = Food101Dataset(
        data_dir=food101_dir,
        transform=data_transforms['train'],
        split='train'
    )
    
    # 划分Food-101训练集，保留一小部分作为验证和测试
    food101_dataset_size = len(food101_train)
    train_size = int(food101_dataset_size * 0.7)
    val_size = int(food101_dataset_size * 0.15)
    test_size = food101_dataset_size - train_size - val_size
    
    food101_train_subset, food101_val, food101_test = torch.utils.data.random_split(
        food101_train, [train_size, val_size, test_size]
    )
    
    # 创建本地数据集
    local_train = CustomFoodDataset(
        root_dir=local_data_dir,
        list_file='train_list.txt',
        name_file='name.txt',
        transform=data_transforms['train']
    )
    
    local_val = CustomFoodDataset(
        root_dir=local_data_dir,
        list_file='val_list.txt',
        name_file='name.txt',
        transform=data_transforms['val']
    )
    
    local_test = CustomFoodDataset(
        root_dir=local_data_dir,
        list_file='test_list.txt',
        name_file='name.txt',
        transform=data_transforms['test']
    )
    
    # 为了保持数据平衡，我们从较大的数据集中采样与较小的数据集相同数量的样本
    # 这里假设local_train比food101_train_subset小
    if len(local_train) < len(food101_train_subset):
        food101_train_indices = torch.randperm(len(food101_train_subset))[:len(local_train)]
        food101_train_subset = torch.utils.data.Subset(food101_train_subset, food101_train_indices)
    elif len(food101_train_subset) < len(local_train):
        local_train_indices = torch.randperm(len(local_train))[:len(food101_train_subset)]
        local_train = torch.utils.data.Subset(local_train, local_train_indices)
    
    # 同样平衡验证和测试集
    if len(local_val) < len(food101_val):
        food101_val_indices = torch.randperm(len(food101_val))[:len(local_val)]
        food101_val = torch.utils.data.Subset(food101_val, food101_val_indices)
    elif len(food101_val) < len(local_val):
        local_val_indices = torch.randperm(len(local_val))[:len(food101_val)]
        local_val = torch.utils.data.Subset(local_val, local_val_indices)
    
    if len(local_test) < len(food101_test):
        food101_test_indices = torch.randperm(len(food101_test))[:len(local_test)]
        food101_test = torch.utils.data.Subset(food101_test, food101_test_indices)
    elif len(food101_test) < len(local_test):
        local_test_indices = torch.randperm(len(local_test))[:len(food101_test)]
        local_test = torch.utils.data.Subset(local_test, local_test_indices)
    
    # 组合数据集
    train_dataset = ConcatDataset([food101_train_subset, local_train])
    val_dataset = ConcatDataset([food101_val, local_val])
    test_dataset = ConcatDataset([food101_test, local_test])
    
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
    
    return dataloaders, dataset_sizes

# 训练二分类器
def train_binary_classifier(model, dataloaders, dataset_sizes, criterion, optimizer, scheduler, num_epochs=5, device='cpu'):
    """训练二分类模型"""
    since = time.time()
    
    best_model_wts = copy.deepcopy(model.state_dict())
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
                running_corrects += torch.sum(preds == labels)
            
            if phase == 'train' and scheduler is not None and not isinstance(scheduler, torch.optim.lr_scheduler.ReduceLROnPlateau):
                scheduler.step()
            
            epoch_loss = running_loss / dataset_sizes[phase]
            epoch_acc = running_corrects.double() / dataset_sizes[phase]
            
            print(f'{phase} Loss: {epoch_loss:.4f} Acc: {epoch_acc:.4f}')
            
            # 如果是验证阶段并且有更好的准确率，保存模型
            if phase == 'val' and epoch_acc > best_acc:
                best_acc = epoch_acc
                best_model_wts = copy.deepcopy(model.state_dict())
                # 保存最佳模型
                torch.save({
                    'model_state_dict': best_model_wts,
                    'optimizer_state_dict': optimizer.state_dict(),
                    'epoch': epoch,
                    'accuracy': best_acc
                }, 'best_binary_classifier.pth')
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

# 更高效的二分类器
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

# 训练二分类器的改进版本（使用更高效的模型）
def main_efficient_binary_classifier():
    # 设置随机种子
    torch.manual_seed(42)
    np.random.seed(42)
    
    # 设置设备
    device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
    print(f"使用设备: {device}")
    
    # 数据集路径
    food101_dir = 'data/food-101'  # 需要修改为真实的Food-101数据集路径
    local_data_dir = 'data/dataset_release/release_data'  # 需要修改为真实的本地数据集路径（此模型用的是ChineseFoodNet）
    
    # 准备数据集
    dataloaders, dataset_sizes = prepare_binary_datasets(food101_dir, local_data_dir)
    print(f"数据集大小: 训练={dataset_sizes['train']}, 验证={dataset_sizes['val']}, 测试={dataset_sizes['test']}")
    
    # 创建二分类模型
    model = EfficientBinaryClassifier(base_model_name='resnet50')
    model = model.to(device)
    
    # 冻结特征提取部分的参数
    for param in model.features.parameters():
        param.requires_grad = False
    
    # 仅训练分类器部分
    optimizer = torch.optim.Adam(model.classifier.parameters(), lr=0.001)
    criterion = nn.CrossEntropyLoss()
    scheduler = torch.optim.lr_scheduler.StepLR(optimizer, step_size=3, gamma=0.1)
    
    # 训练模型的第一阶段（只训练分类器）
    print("训练第一阶段（只训练分类器）...")
    model = train_binary_classifier(
        model=model,
        dataloaders=dataloaders,
        dataset_sizes=dataset_sizes,
        criterion=criterion,
        optimizer=optimizer,
        scheduler=scheduler,
        num_epochs=5,
        device=device
    )
    
    # 解冻特征提取部分的最后一层
    print("训练第二阶段（微调特征提取器的最后一层）...")
    for param in model.features[-1].parameters():
        param.requires_grad = True
    
    # 使用较小的学习率训练整个模型
    optimizer = torch.optim.Adam([
        {'params': model.classifier.parameters(), 'lr': 0.0001},
        {'params': model.features[-1].parameters(), 'lr': 0.00001}
    ])
    scheduler = torch.optim.lr_scheduler.StepLR(optimizer, step_size=3, gamma=0.1)
    
    # 训练模型的第二阶段
    model = train_binary_classifier(
        model=model,
        dataloaders=dataloaders,
        dataset_sizes=dataset_sizes,
        criterion=criterion,
        optimizer=optimizer,
        scheduler=scheduler,
        num_epochs=5,
        device=device
    )
    
    # 保存完整模型
    torch.save({
        'model_state_dict': model.state_dict()
    }, 'efficient_binary_classifier.pth')
    print("完整模型已保存为 'efficient_binary_classifier.pth'")
    
    # 评估最终模型
    model.eval()
    running_corrects = 0
    class_correct = [0, 0]
    class_total = [0, 0]
    
    with torch.no_grad():
        for inputs, labels in dataloaders['test']:
            inputs = inputs.to(device)
            labels = labels.to(device)
            
            outputs = model(inputs)
            _, preds = torch.max(outputs, 1)
            
            running_corrects += torch.sum(preds == labels)
            
            for i in range(len(labels)):
                label = labels[i].item()
                class_correct[label] += (preds[i] == label).item()
                class_total[label] += 1
    # todo: 测试集
    # test_acc = running_corrects / dataset_sizes['test']
    # print(f'测试准确率: {test_acc:.4f}')
    
    for i in range(2):
        class_name = "Food-101" if i == 0 else "本地数据库"
        if class_total[i] > 0:
            print(f'{class_name} 准确率: {class_correct[i] / class_total[i]:.4f} ({class_correct[i]}/{class_total[i]})')
            
if __name__ == "__main__":
    main_efficient_binary_classifier()
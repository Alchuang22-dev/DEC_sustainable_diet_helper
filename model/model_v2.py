import torch
import torch.nn as nn
import torch.optim as optim
from torch.optim import lr_scheduler
from torch.utils.data import DataLoader
import torchvision
from torchvision import models, transforms, datasets
import time
import copy
import os

# 定义数据转换（重要：要与预训练模型使用的规范化参数一致）
data_transforms = {
    'train': transforms.Compose([
        transforms.RandomResizedCrop(224),  # 随机裁剪到224x224
        transforms.RandomHorizontalFlip(),   # 随机水平翻转
        transforms.ColorJitter(brightness=0.2, contrast=0.2, saturation=0.2), # 色彩增强
        transforms.ToTensor(),
        transforms.Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])  # ImageNet标准化参数
    ]),
    'test': transforms.Compose([
        transforms.Resize(256),              # 调整大小
        transforms.CenterCrop(224),          # 中心裁剪
        transforms.ToTensor(),
        transforms.Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])
    ]),
}

# 使用您之前修改的Food101数据集类
data_dir = './data'  # 数据集根目录
train_dataset = datasets.Food101(root=data_dir, split='train', transform=data_transforms['train'])
test_dataset = datasets.Food101(root=data_dir, split='test', transform=data_transforms['test'])

# 创建数据加载器
train_loader = DataLoader(train_dataset, batch_size=32, shuffle=True, num_workers=4)
test_loader = DataLoader(test_dataset, batch_size=32, shuffle=False, num_workers=4)

dataloaders = {'train': train_loader, 'test': test_loader}
dataset_sizes = {'train': len(train_dataset), 'test': len(test_dataset)}
class_names = train_dataset.classes

def create_model(model_name='resnet18', pretrained=True, freeze_backbone=False):
    """
    创建并配置预训练模型
    
    参数:
    - model_name: 使用的模型架构（'resnet18', 'resnet50', 'resnet101'等）
    - pretrained: 是否使用预训练权重
    - freeze_backbone: 是否冻结主干网络参数
    """
    # 加载预训练模型
    if model_name == 'resnet18':
        model = models.resnet18(weights='IMAGENET1K_V1' if pretrained else None)
    elif model_name == 'resnet50':
        model = models.resnet50(weights='IMAGENET1K_V1' if pretrained else None)
    elif model_name == 'resnet101':
        model = models.resnet101(weights='IMAGENET1K_V1' if pretrained else None)
    else:
        raise ValueError(f"不支持的模型: {model_name}")
    
    # 如果需要冻结主干网络
    if freeze_backbone:
        for param in model.parameters():
            param.requires_grad = False
    
    # 修改最后的全连接层以适应Food-101的101个类别
    num_features = model.fc.in_features
    model.fc = nn.Sequential(
        nn.Dropout(0.5),
        nn.Linear(num_features, 101)
    )
    
    return model


def train_model(model, criterion, optimizer, scheduler, num_epochs=25):
    """
    训练模型
    
    参数:
    - model: 需要训练的模型
    - criterion: 损失函数
    - optimizer: 优化器
    - scheduler: 学习率调度器
    - num_epochs: 训练的轮数
    """
    since = time.time()

    best_model_wts = copy.deepcopy(model.state_dict())
    best_acc = 0.0

    for epoch in range(num_epochs):
        print(f'Epoch {epoch+1}/{num_epochs}')
        print('-' * 10)

        # 每个epoch有训练和验证阶段
        for phase in ['train', 'test']:
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
                # 只有在训练时跟踪梯度
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

            epoch_loss = running_loss / dataset_sizes[phase]
            epoch_acc = running_corrects.double() / dataset_sizes[phase]

            print(f'{phase} Loss: {epoch_loss:.4f} Acc: {epoch_acc:.4f}')

            # 保存最佳模型（基于测试集精度）
            if phase == 'test' and epoch_acc > best_acc:
                best_acc = epoch_acc
                best_model_wts = copy.deepcopy(model.state_dict())
                # 保存最佳模型
                torch.save(best_model_wts, 'best_model_food101.pth')
                print(f'保存新的最佳模型，精度: {best_acc:.4f}')
            
            # 在验证阶段结束后，如果使用ReduceLROnPlateau调度器，则传入验证损失
            if phase == 'test' and scheduler is not None and isinstance(scheduler, torch.optim.lr_scheduler.ReduceLROnPlateau):
                scheduler.step(epoch_loss)  # 使用验证集的损失来调整学习率

        # 如果使用其他类型的调度器（非ReduceLROnPlateau），在每个epoch结束后调用step()
        if scheduler is not None and not isinstance(scheduler, torch.optim.lr_scheduler.ReduceLROnPlateau):
            scheduler.step()

        print()

    time_elapsed = time.time() - since
    print(f'训练完成，用时 {time_elapsed // 60:.0f}m {time_elapsed % 60:.0f}s')
    print(f'最佳测试精度: {best_acc:.4f}')

    # 加载最佳模型权重
    model.load_state_dict(best_model_wts)
    return model

# 创建模型，冻结主干网络
model = create_model('resnet50', pretrained=True, freeze_backbone=True)
device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
model = model.to(device)

# 定义损失函数
criterion = nn.CrossEntropyLoss()

# 只优化最后一层（全连接层）的参数
optimizer = optim.SGD(model.fc.parameters(), lr=0.001, momentum=0.9)

# 学习率调度器：每7个epoch学习率降低10倍
scheduler = lr_scheduler.StepLR(optimizer, step_size=7, gamma=0.1)

# 训练模型 - 先训练10个epoch
model = train_model(model, criterion, optimizer, scheduler, num_epochs=10)

# 解冻所有层
for param in model.parameters():
    param.requires_grad = True

# 使用较小的学习率优化所有参数
# 对不同的层组使用不同的学习率
# 一般原则：靠近输入的层使用较小的学习率，靠近输出的层使用较大的学习率

# 将模型分成不同的参数组
params_to_update = []
# 较低学习率用于卷积层
params_to_update.append({'params': model.conv1.parameters(), 'lr': 0.0001})
params_to_update.append({'params': model.bn1.parameters(), 'lr': 0.0001})
params_to_update.append({'params': model.layer1.parameters(), 'lr': 0.0001})
params_to_update.append({'params': model.layer2.parameters(), 'lr': 0.0001})
# 中等学习率用于中间层
params_to_update.append({'params': model.layer3.parameters(), 'lr': 0.0005})
params_to_update.append({'params': model.layer4.parameters(), 'lr': 0.0005})
# 较高学习率用于分类器
params_to_update.append({'params': model.fc.parameters(), 'lr': 0.001})

# 使用Adam优化器进行细粒度微调
optimizer = optim.Adam(params_to_update)

# 学习率调度器
scheduler = lr_scheduler.ReduceLROnPlateau(optimizer, mode='min', factor=0.1, patience=2, verbose=True)

# 继续训练模型，进行完整微调
model = train_model(model, criterion, optimizer, scheduler, num_epochs=15)

def evaluate_model(model, dataloader):
    model.eval()
    running_corrects = 0
    total = 0
    
    # 禁用梯度计算
    with torch.no_grad():
        for inputs, labels in dataloader:
            inputs = inputs.to(device)
            labels = labels.to(device)
            
            # 前向传播
            outputs = model(inputs)
            _, preds = torch.max(outputs, 1)
            
            # 计算正确预测的数量
            running_corrects += torch.sum(preds == labels.data)
            total += labels.size(0)
    
    # 计算准确率
    accuracy = running_corrects.double() / total
    print(f'测试集准确率: {accuracy:.4f}')
    return accuracy

# 评估最终模型
final_accuracy = evaluate_model(model, dataloaders['test'])

# 保存完整模型
torch.save({
    'model_state_dict': model.state_dict(),
    'optimizer_state_dict': optimizer.state_dict(),
    'class_names': class_names
}, 'food101_model_complete.pth')


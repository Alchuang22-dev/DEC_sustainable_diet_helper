import torch
from torchvision import datasets, transforms
import matplotlib.pyplot as plt

def show_dataset():
    # 数据预处理
    transform = transforms.Compose([
        transforms.Resize((256, 256)),  # 调整图片大小便于显示
        transforms.ToTensor(),
        transforms.Normalize((0.5, 0.5, 0.5), (0.5, 0.5, 0.5))
    ])

    # 加载训练集
    print("Loading Food-101 dataset...")
    trainset = datasets.Food101(root='./data', split='train', download=True, transform=transform)
    
    # 获取类别信息
    classes = trainset.classes
    
    # 显示数据集信息
    print(f"\nDataset size: {len(trainset)} images")
    sample_image, sample_label = trainset[0]
    print(f"Image shape: {sample_image.shape}")
    print(f"Number of classes: {len(classes)}")
    
    # 显示样本图片
    plt.figure(figsize=(15, 5))
    for i in range(5):
        img, label = trainset[i]
        # 反归一化：原归一化公式是 img = (img-mean)/std, 此处逆操作为 img*std+mean
        img = img * 0.5 + 0.5
        img = img.numpy()
        plt.subplot(1, 5, i + 1)
        plt.imshow(img.transpose(1, 2, 0))
        plt.title(f'Class: {classes[label]}')
        plt.axis('off')
    plt.tight_layout()
    plt.savefig('food101_samples.png')
    
    # 显示每个类别的样本数量
    class_counts = torch.zeros(len(classes))
    for _, label in trainset:
        class_counts[label] += 1
    
    plt.figure(figsize=(15, 8))
    plt.bar(classes, class_counts)
    plt.title('Number of samples per class')
    plt.xticks(rotation=90)
    plt.tight_layout()
    plt.savefig('food101_class_distribution.png')
    
    print("\nVisualization files saved:")
    print("- food101_samples.png")
    print("- food101_class_distribution.png")

if __name__ == '__main__':
    show_dataset()

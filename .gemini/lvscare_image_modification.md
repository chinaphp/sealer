# LvsCare 镜像地址修改说明

## 修改目的
解决 LvsCare 镜像使用华为云在线地址在离线/网络受限环境下无法部署的问题，并使用正确的镜像版本 `sealerio/lvscare:v1.1.3-beta.8`。

## 修改内容

### 0. 镜像版本常量定义
**文件**: `common/common.go`
**行号**: 49
**修改前**:
```go
LvsCareRepoAndTag = "docker-cdv5ju.swr-pro.myhuaweicloud.com/global/lvscare:latest"
```

**修改后**:
```go
LvsCareRepoAndTag = "sealerio/lvscare:v1.1.3-beta.8"
```

**说明**: 
- 移除了华为云在线地址前缀
- 使用实际打包的镜像名称和版本号
- 便于集中管理和升级镜像版本

### 1. Registry 负载均衡场景
**文件**: `pkg/registry/local.go`
**行号**: 182-184
**修改前**:
```go
//todo should make lvs image name as const value in sealer repo.
lvsImageURL := common.LvsCareRepoAndTag  // "docker-cdv5ju.swr-pro.myhuaweicloud.com/global/lvscare:latest"
```

**修改后**:
```go
// Use local registry for lvscare image to support offline deployment
// Format: registry.domain:port/sealerio/lvscare:v1.1.3-beta.8
lvsImageURL := fmt.Sprintf("%s/%s", net.JoinHostPort(c.Domain, strconv.Itoa(c.Port)), common.LvsCareRepoAndTag)
```

**说明**: 
- 在配置 Registry 高可用的 IPVS 负载均衡时，使用本地 registry 的域名和端口构建 lvscare 镜像地址
- 引用 `common.LvsCareRepoAndTag` 常量，确保版本统一管理
- 例如：`sea.hub:5000/sealerio/lvscare:v1.1.3-beta.8`

### 2. Kubernetes API Server 负载均衡场景
**文件**: `pkg/runtime/kubernetes/utils.go`
**行号**: 190-193
**修改前**:
```go
lvsImageURL := common.LvsCareRepoAndTag  // "docker-cdv5ju.swr-pro.myhuaweicloud.com/global/lvscare:latest"
```

**修改后**:
```go
// Use local registry for lvscare image to support offline deployment
// Get registry URL from cluster config and construct image reference
// Format: registry.url/sealerio/lvscare:v1.1.3-beta.8
lvsImageURL := fmt.Sprintf("%s/%s", k.Config.RegistryInfo.URL, common.LvsCareRepoAndTag)
```

**说明**:
- 在配置 Kubernetes Master 高可用的 IPVS 负载均衡时，从集群配置中获取 registry URL
- 引用 `common.LvsCareRepoAndTag` 常量，确保版本统一管理
- 动态构建镜像地址，例如：`sea.hub:5000/sealerio/lvscare:v1.1.3-beta.8`

## 技术细节

### 镜像拉取策略
静态 Pod 配置中已设置 `ImagePullPolicy: IfNotPresent`，这意味着：
- 如果本地已有镜像，kubelet 不会尝试拉取
- 如果本地没有镜像，才会从指定的 registry 拉取

### 镜像预加载要求
为确保离线部署成功，需要：
1. **lvscare 镜像打包到 ClusterImage**: 确保在构建 ClusterImage 时包含 `sealerio/lvscare:v1.1.3-beta.8` 镜像
2. **镜像命名统一**: 在本地 registry 中，镜像完整路径应为 `sea.hub:5000/sealerio/lvscare:v1.1.3-beta.8`
3. **版本管理**: 通过修改 `common.LvsCareRepoAndTag` 常量统一升级版本

## 影响范围

### 正面影响
✅ **支持离线部署**: 不再依赖华为云在线镜像源
✅ **网络受限环境**: 在防火墙内或无法访问公网的环境下可以正常部署
✅ **统一镜像管理**: 所有镜像都从本地 registry 获取，便于版本控制
✅ **加快部署速度**: 从本地 registry 拉取镜像比从公网拉取更快
✅ **版本固定**: 使用明确的版本号 `v1.1.3-beta.8`，避免 `latest` 标签的不确定性
✅ **集中管理**: 通过 `common.LvsCareRepoAndTag` 常量统一管理版本，便于升级

### 潜在影响
⚠️ **镜像预加载依赖**: 必须确保 `sealerio/lvscare:v1.1.3-beta.8` 镜像已经加载到本地 registry
⚠️ **构建流程修改**: 需要检查 ClusterImage 构建流程是否包含正确版本的 lvscare 镜像
⚠️ **版本升级**: 升级 lvscare 版本时需要同步修改 `common.LvsCareRepoAndTag` 常量

## 测试建议

1. **离线环境测试**:
   ```bash
   # 确保本地 registry 包含 lvscare 镜像
   nerdctl images | grep lvscare
   # 应显示：sea.hub:5000/sealerio/lvscare   v1.1.3-beta.8
   
   # 部署高可用集群
   sealer run -f Clusterfile
   
   # 检查静态 Pod 状态
   kubectl get pods -n kube-system | grep lvscare
   ```

2. **验证镜像地址**:
   ```bash
   # 检查静态 Pod 配置
   cat /etc/kubernetes/manifests/kube-lvscare.yaml | grep image:
   cat /etc/kubernetes/manifests/reg-lvscare.yaml | grep image:
   
   # 应显示本地 registry 地址，例如：
   # image: sea.hub:5000/sealerio/lvscare:v1.1.3-beta.8
   ```

3. **验证镜像可用性**:
   ```bash
   # 检查镜像是否已缓存到节点
   nerdctl -n k8s.io images | grep lvscare
   
   # 测试从本地 registry 拉取镜像
   nerdctl pull sea.hub:5000/sealerio/lvscare:v1.1.3-beta.8
   ```

## 后续优化建议

1. ✅ **使用语义化版本**（已实现）: 已将 `lvscare:latest` 改为固定版本号 `sealerio/lvscare:v1.1.3-beta.8`
2. ✅ **集中管理版本**（已实现）: 通过 `common.LvsCareRepoAndTag` 常量统一管理版本
3. **镜像存在性检查**: 在部署前检查 lvscare 镜像是否存在于本地 registry
4. **降级策略**: 如果本地 registry 不可用，考虑使用节点本地已缓存的镜像
5. **配置化增强**: 允许用户通过环境变量覆盖 lvscare 镜像地址（用于特殊场景）

## 兼容性说明

- ✅ 向后兼容：使用本地 registry 不影响现有功能
- ✅ 支持多架构：镜像地址构建逻辑与架构无关
- ✅ 支持 IPv6：镜像地址使用 `net.JoinHostPort` 处理 IPv6 地址

## 修改日期
2026-01-11

## 相关 Issue
- 离线环境部署失败：无法从华为云拉取 lvscare 镜像
- 网络受限环境下集群高可用功能无法使用

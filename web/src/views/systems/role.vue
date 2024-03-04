<template>
  <div class="main-container">
    <TableBody>
      <template #header>
        <TableHeader :show-filter="false">
          <template #table-config>
            <AddButton @add="onAddItem"/>
          </template>
        </TableHeader>
      </template>
      <template #default>
        <a-table
            :bordered="{ wrapper: true, cell: true }"
            :data="dataList"
            :loading="tableLoading"
            :pagination="false"
            :row-key="rowKey"
        >
          <template #columns>
            <a-table-column
                v-for="item of tableColumns"
                :key="item.key"
                :align="item.align"
                :data-index="(item.key as string)"
                :fixed="item.fixed"
                :title="(item.title as string)"
                :width="item.width"
            >
              <template v-if="item.key === 'index'" #cell="{ rowIndex }">
                {{ rowIndex + 1 }}
              </template>
              <template v-else-if="item.key === 'gender'" #cell="{ record }">
                {{ record.gender === 1 ? '男' : '女' }}
              </template>
              <template v-else-if="item.key === 'avatar'" #cell="{ record }">
                <a-avatar>
                  <img :src="record.avatar"/>
                </a-avatar>
              </template>
              <template v-else-if="item.key === 'actions'" #cell="{ record }">
                <a-space>
                  <a-button size="mini" status="success" @click="onUpdateItem(record)"
                  >编辑
                  </a-button
                  >
                  <a-button size="mini" status="danger" @click="onDeleteItem(record)"
                  >删除
                  </a-button
                  >
                  <a-button size="mini" status="warning" @click="onShowMenu(record)">
                    菜单权限
                  </a-button>
                </a-space>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </template>
    </TableBody>
    <ModalDialog ref="modalDialogRef" :title="actionTitle" @confirm="onDataFormConfirm">
      <template #content>
        <a-form :model="formModel">
          <a-form-item
              v-for="item of formItems"
              :key="item.key"
              :class="[item.required ? 'form-item__require' : 'form-item__no_require']"
              :label="item.label"
          >
            <template v-if="item.type === 'input'">
              <a-input v-model="item.value.value" :placeholder="item.placeholder">
                <template v-if="item.key === 'roleCode'" #prepend>
                  {{ 123 }}
                </template>
              </a-input>
            </template>
            <template v-if="item.type === 'textarea'">
              <a-textarea
                  v-model="item.value.value"
                  :auto-size="{ minRows: 3, maxRows: 5 }"
                  :placeholder="item.placeholder"
              />
            </template>
          </a-form-item>
        </a-form>
      </template>
    </ModalDialog>
    <ModalDialog ref="menuModalDialogRef" title="编辑菜单权限" @confirm="onDataFormConfirm">
      <template #content>
        <a-tree
            v-model:checked-keys="defaultCheckedKeys"
            v-model:expanded-keys="defaultExpandedKeys"
            :data="menuData"
            checkable
        />
      </template>
    </ModalDialog>
  </div>
</template>

<script lang="ts">
import {post} from '@/api/http'
import {getMenuListByRoleId, getRoleList} from '@/api/url'
import {useRowKey, useTable, useTableColumn} from '@/hooks/table'
import {FormItem, ModalDialogType} from '@/types/components'
import {Message, Modal} from '@arco-design/web-vue'
import {defineComponent, nextTick, onMounted, ref} from 'vue'
import useUserStore from "@/store/modules/user";

const ROLE_CODE_FLAG = 'ROLE_'
const formItems = [
  {
    label: '角色名称',
    type: 'input',
    key: 'name',
    value: ref(''),

    required: true,
    placeholder: '请输入角色名称',
    validator: function () {
      if (!this.value.value) {
        Message.error(this.placeholder || '')
        return false
      }
      return true
    },
  },
  {
    label: '角色编号',
    key: 'roleCode',
    value: ref(''),
    type: 'input',
    required: true,
    placeholder: '请输入角色编号',
    validator: function () {
      if (!this.value.value) {
        Message.error(this.placeholder || '')
        return false
      }
      return true
    },
  },
  {
    label: '角色描述',
    key: 'description',
    value: ref(''),
    type: 'textarea',
    placeholder: '请输入角色描述',
  },
] as FormItem[]

function handleMenuData(
    menuData: Array<any>,
    defaultCheckedKeys: Array<string>,
    defaultExpandedKeys: Array<string>
) {
  const tempMenus = [] as Array<any>
  menuData.forEach((it) => {
    const tempMenu = {} as any
    tempMenu.key = it.menuUrl
    tempMenu.title = it.menuName
    defaultCheckedKeys.push(tempMenu.key as string)
    if (it.children) {
      defaultExpandedKeys.push(tempMenu.key as string)
      tempMenu.children = handleMenuData(it.children, defaultCheckedKeys, defaultExpandedKeys)
    }
    tempMenus.push(tempMenu)
  })
  return tempMenus
}

export default defineComponent({
  name: 'Role',
  setup() {
    const modalDialogRef = ref<ModalDialogType | null>(null)
    const menuModalDialogRef = ref<ModalDialogType | null>(null)
    const table = useTable()
    const rowKey = useRowKey('id')
    const actionTitle = ref('添加角色')
    const menuData = ref([] as Array<any>)
    const tableColumns = useTableColumn([
      table.indexColumn,
      {
        title: '角色名称',
        key: 'name',
        dataIndex: 'name',
      },
      {
        title: '角色编号',
        key: 'roleCode',
        dataIndex: 'roleCode',
      },
      {
        title: '角色描述',
        key: 'description',
        dataIndex: 'description',
      },
      {
        title: '创建时间',
        key: 'createTime',
        dataIndex: 'createTime',
      },
      {
        title: '操作',
        key: 'actions',
        dataIndex: 'actions',
      },
    ])
    const defaultCheckedKeys = ref([] as Array<string>)
    const defaultExpandedKeys = ref([] as Array<string>)
    const userStore = useUserStore()
    const formModel = ref({})
    let num = 1

    function doRefresh() {
      post({
        url: getRoleList,
        data: {},
      })
          .then(table.handleSuccess)
          .catch(console.log)
    }

    function onAddItem() {
      num = 1
      actionTitle.value = '添加角色'
      modalDialogRef.value?.toggle()
      formItems.forEach((it) => {
        if (it.reset) {
          it.reset()
        } else {
          it.value.value = ''
        }
      })
    }

    function onUpdateItem(item: any) {
      actionTitle.value = '编辑角色'
      num = 2
      modalDialogRef.value?.toggle()
      nextTick(() => {
        formItems.forEach((it) => {
          const key = it.key
          const propName = item[key]
          if (propName) {
            if (it.key === 'roleCode') {
              it.value.value = propName.replace(ROLE_CODE_FLAG, '')
            } else {
              it.value.value = propName
            }
          }
        })
      })
    }

    function onDeleteItem(data: any) {
      Modal.confirm({
        title: '提示',
        content: '是否要删除此角色？',
        cancelText: '取消',
        okText: '删除',
        onOk: () => {
          Message.success('模拟角色删除成功，参数为' + JSON.stringify(data))
        },
      })
    }

    function onDataFormConfirm() {
      if (num === 1) {

      }
      if (num == 2) {
      }
      if (num == 3) {
      }
      menuModalDialogRef.value?.toggle()

      // if (formItems.every((it) => (it.validator ? it.validator() : true))) {
      //   modalDialogRef.value?.toggle()
      //   Message.success(
      //       '模拟菜单添加成功，参数为：' +
      //       JSON.stringify(
      //           formItems.reduce((pre, cur) => {
      //             ;(pre as any)[cur.key] = cur.value.value
      //             return pre
      //           }, {})
      //       )
      //   )
      // }
    }

    function onShowMenu(item: any) {
      num = 3
      post({
        url: getMenuListByRoleId,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
        data: {
          roleId: item.id,
        },
      })
          .then((res) => {
            menuData.value = []
            defaultCheckedKeys.value = []
            defaultExpandedKeys.value = []
            menuData.value = handleMenuData(
                res.data,
                defaultCheckedKeys.value,
                defaultExpandedKeys.value
            )
            menuModalDialogRef.value?.toggle()
          })
          .catch(console.log)
    }

    onMounted(doRefresh)
    return {
      ROLE_CODE_FLAG,
      modalDialogRef,
      menuModalDialogRef,
      rowKey,
      formModel,
      menuData,
      tableColumns,
      formItems,
      actionTitle,
      defaultCheckedKeys,
      defaultExpandedKeys,
      ...table,
      num,
      onAddItem,
      onDataFormConfirm,
      onShowMenu,
      onDeleteItem,
      onUpdateItem,
    }
  },
})
</script>

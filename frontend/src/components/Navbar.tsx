import {
    AppShell,
    Avatar,
    Badge,
    Box,
    Button,
    Group,
    Progress,
    Stack,
    Text,
    UnstyledButton,
} from '@mantine/core';
import {
  IconChartBar,
  IconCategory,
  IconFileText,
  IconLayoutDashboard,
  IconList,
  IconUpload,
} from '@tabler/icons-react';
import classes from './Navbar.module.css';

export type NavbarItemKey = 
    | 'overview'
    | 'statements'
    | 'transactions'
    | 'categories'
    | 'insights';

export interface NavbarProps {
    current: NavbarItemKey;
    onNavigate: (key: NavbarItemKey) => void;
    onUpload: () => void;
    monthStatus: {
        label: string;
        uploaded: number;
        expected: number;
        detail: string;
    };
    statementsMissing?: number;
    user: {
        name: string;
        email: string;
        initial?: string;
    }
}

interface NavItem {
    key: NavbarItemKey;
    label: string;
    icon: typeof IconLayoutDashboard;
    badge?: number;
}

export function Navbar({
  current,
  onNavigate,
  onUpload,
  monthStatus,
  statementsMissing,
  user,
}: NavbarProps) {
  const items: NavItem[] = [
    { key: 'overview', label: 'Overview', icon: IconLayoutDashboard },
    {
      key: 'statements',
      label: 'Statements',
      icon: IconFileText,
      badge: statementsMissing,
    },
    { key: 'transactions', label: 'Transactions', icon: IconList },
    { key: 'categories', label: 'Categories', icon: IconCategory },
    { key: 'insights', label: 'Insights', icon: IconChartBar },
  ];
 
  const progress = monthStatus.expected
    ? (monthStatus.uploaded / monthStatus.expected) * 100
    : 0;
 
  return (
    <AppShell.Navbar p="md" className={classes.navbar}>
      {/* Brand */}
      <Group gap="xs" px="xs" pb="xl" className={classes.brand}>
        <Box className={classes.brandMark}>K</Box>
        <Text fw={600} size="md">
          Kinji
        </Text>
      </Group>
 
      {/* Primary action */}
      <Button
        fullWidth
        leftSection={<IconUpload size={16} />}
        onClick={onUpload}
        className={classes.uploadButton}
        mb="lg"
      >
        Upload statement
      </Button>
 
      {/* Section label */}
      <Text className={classes.sectionLabel} px="xs" pb={6}>
        Workspace
      </Text>
 
      {/* Nav items */}
      <Stack gap={1}>
        {items.map((item) => {
          const Icon = item.icon;
          const isActive = item.key === current;
          return (
            <UnstyledButton
              key={item.key}
              onClick={() => onNavigate(item.key)}
              className={classes.navItem}
              data-active={isActive || undefined}
            >
              <Icon size={16} className={classes.navIcon} stroke={1.75} />
              <Text size="sm" className={classes.navLabel}>
                {item.label}
              </Text>
              {item.badge ? (
                <Badge
                  size="sm"
                  variant="light"
                  className={classes.navBadge}
                  ml="auto"
                >
                  {item.badge} missing
                </Badge>
              ) : null}
            </UnstyledButton>
          );
        })}
      </Stack>
 
      {/* Spacer pushes the rest to the bottom */}
      <Box style={{ flex: 1 }} />
 
      {/* Month status */}
      <Box className={classes.statusBlock} px="xs" py="md">
        <Text className={classes.statusLabel} mb={8}>
          {monthStatus.label}
        </Text>
        <Text size="sm" fw={500} mb={6}>
          {monthStatus.uploaded} of {monthStatus.expected} statements
        </Text>
        <Progress
          value={progress}
          size="xs"
          color="dark"
          radius="xl"
          className={classes.statusProgress}
          mb={10}
        />
        <Text size="xs" c="dimmed" lh={1.4}>
          {monthStatus.detail}
        </Text>
      </Box>
 
      {/* User chip */}
      <UnstyledButton className={classes.userChip} mt={8}>
        <Group gap="xs" wrap="nowrap">
          <Avatar size={28} radius="xl" color="gray">
            {user.initial ?? user.name.charAt(0)}
          </Avatar>
          <Box style={{ minWidth: 0, flex: 1 }}>
            <Text size="sm" fw={500} truncate>
              {user.name}
            </Text>
            <Text size="xs" c="dimmed" truncate>
              {user.email}
            </Text>
          </Box>
        </Group>
      </UnstyledButton>
    </AppShell.Navbar>
  );
}